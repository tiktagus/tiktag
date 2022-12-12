/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codenotary/immudb/embedded/sql"
	"github.com/codenotary/immudb/embedded/store"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sony/sonyflake"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func fakeMachineID(uint16) bool {
	return true
}

type Product struct {
	ID     int `gorm:"primarykey"`
	Code   string
	Price  uint
	Amount uint
}

func GetFileContentType(ouput *os.File) (string, error) {
	// to sniff the content type only the first
	// 512 bytes are used.
	buf := make([]byte, 512)
	_, err := ouput.Read(buf)
	if err != nil {
		return "", err
	}
	// the function that actually does the trick
	contentType := http.DetectContentType(buf)
	return contentType, nil
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tiktag [file to upload]",
	Short: "A command-line tool for preparing images for blog post or sharing.",
	Long:  `Upload a photo and get its S3 URL back as a response, for use in Markdown for publishing.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// sha256
		fn := args[0]
		f, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		contentType, err := GetFileContentType(f)
		if err != nil {
			log.Fatal(err)
		}

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%x\n", h.Sum(nil))

		// Sonyflake Id
		var st sonyflake.Settings
		st.CheckMachineID = fakeMachineID
		sf := sonyflake.NewSonyflake(st)
		if sf == nil {
			log.Fatal("New Sonyflake failed!")
		}

		id, err := sf.NextID()
		if err != nil {
			log.Fatal("NextID failed!")
		}

		fmt.Println(id)

		// Minio
		ctx := context.Background()
		endpoint := viper.GetString("minio.endpoint")
		accessKeyID := viper.GetString("minio.id")
		secretAccessKey := viper.GetString("minio.secret")
		useSSL := viper.GetBool("minio.useSSL")

		// Initialize minio client object.
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}

		// Make a new bucket called mymusic.
		bucketName := viper.GetString("minio.bucketName")

		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
			if errBucketExists == nil && exists {
				log.Printf("We already own %s\n", bucketName)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Printf("Successfully created %s\n", bucketName)
		}

		// Upload the zip file
		_, file := filepath.Split(fn)
		objectName := file
		filePath := fn

		// Upload the zip file with FPutObject
		info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
		log.Println(info)

		// immudb embedded
		// create/open immudb store at specified path
		db, err := store.Open("data", store.DefaultOptions())
		handleErr(err)
		defer db.Close()

		// initialize sql engine (specify a key-prefix to isolate generated kv entries)
		engine, err := sql.NewEngine(db, sql.DefaultOptions().WithPrefix([]byte("sql")))
		handleErr(err)

		_, _, err = engine.Exec("CREATE DATABASE db1;", nil, nil)
		handleErr(err)

		// set the database to use in the context of the ongoing sql tx
		_, _, err = engine.Exec("USE DATABASE db1;", nil, nil)
		handleErr(err)

		// a sql tx is created and carried over next statements
		sqltx, _, err := engine.Exec("BEGIN TRANSACTION;", nil, nil)
		handleErr(err)

		// ensure tx is closed (it won't affect committed tx)
		defer engine.Exec("ROLLBACK;", nil, sqltx)

		// creates a table
		_, _, err = engine.Exec(`
		CREATE TABLE journal (
			id INTEGER,
			date TIMESTAMP,
			creditaccount INTEGER,
			debitaccount INTEGER,
			amount INTEGER,
			description VARCHAR,
			PRIMARY KEY id
		);`, nil, sqltx)
		handleErr(err)

		// insert some rows
		_, _, err = engine.Exec(`
		INSERT INTO journal (
			id,
			date,
			creditaccount,
			debitaccount,
			amount,
			description
		) VALUES 
			(1, NOW(), 100, 0, 4000, 'CREDIT'),
			(2, NOW(), 0, 50, 4100, 'DEBIT')
		;`, nil, sqltx)
		handleErr(err)

		// query data including ongoing and unconfirmed changes
		rowReader, err := engine.Query(`
			SELECT id, date, creditaccount, debitaccount, amount, description
			FROM journal
			WHERE amount > @value;
	`, map[string]interface{}{"value": 100}, sqltx)
		handleErr(err)

		// ensure row reader is closed
		defer rowReader.Close()

		// selected columns can be read from the rowReader
		cols, err := rowReader.Columns()
		handleErr(err)

		for {
			// iterate over result set
			row, err := rowReader.Read()
			if err == sql.ErrNoMoreRows {
				break
			}
			handleErr(err)

			// each row contains values for the selected columns
			log.Printf("row: %v\n", row.ValuesBySelector[cols[0].Selector()].Value())
		}

		// close row reader
		rowReader.Close()

		// commit ongoing transaction
		_, _, err = engine.Exec("COMMIT;", nil, sqltx)
		handleErr(err)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tiktag.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/tiktag/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.tiktag") // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.ReadInConfig()
}
