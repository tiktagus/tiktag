/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/codenotary/immudb/pkg/client"
	// immugorm "github.com/codenotary/immugorm"
	immudb "github.com/0ctanium/gorm-immudb"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
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
		// f, err := os.Open(args[0])
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer f.Close()

		// h := sha256.New()
		// if _, err := io.Copy(h, f); err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Printf("%x\n", h.Sum(nil))

		// Sonyflake Id
		// var st sonyflake.Settings
		// st.CheckMachineID = fakeMachineID
		// sf := sonyflake.NewSonyflake(st)
		// if sf == nil {
		// 	log.Fatal("New Sonyflake failed!")
		// }

		// id, err := sf.NextID()
		// if err != nil {
		// 	log.Fatal("NextID failed!")
		// }

		// fmt.Println(id)

		// Minio
		// ctx := context.Background()
		// endpoint := "s3.tikoly.com"
		// // FIXME: DO NOT commit!!!!
		// accessKeyID := "ZaRTBCf2g4ZMGVgu"
		// secretAccessKey := "xCHRw9vv1etAUx0pABvRQiDLnF4euowj"
		// useSSL := true

		// // Initialize minio client object.
		// minioClient, err := minio.New(endpoint, &minio.Options{
		// 	Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		// 	Secure: useSSL,
		// })
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		// // Make a new bucket called mymusic.
		// bucketName := "tiktag"

		// err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		// if err != nil {
		// 	// Check to see if we already own this bucket (which happens if you run this twice)
		// 	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		// 	if errBucketExists == nil && exists {
		// 		log.Printf("We already own %s\n", bucketName)
		// 	} else {
		// 		log.Fatalln(err)
		// 	}
		// } else {
		// 	log.Printf("Successfully created %s\n", bucketName)
		// }

		// // Upload the zip file
		// objectName := "demo.jpg"
		// filePath := "./demo.jpg"
		// contentType := "image/jpeg"

		// // Upload the zip file with FPutObject
		// info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		// log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

		// immudb
		opts := client.DefaultOptions()

		opts.Username = "immudb"
		opts.Password = "immudb"
		opts.Database = "defaultdb"
		opts.HealthCheckRetries = 10

		// db, err := gorm.Open(immugorm.OpenWithOptions(opts, &immugorm.ImmuGormConfig{Verify: false}), &gorm.Config{
		// 	Logger: logger.Default.LogMode(logger.Info),
		// })
		db, err := gorm.Open(immudb.New(immudb.Config{
			DSN:                "immudb://immudb:immudb@127.0.0.1:3322/defaultdb", // data source name, refer https://docs.immudb.io/master/develop/sqlstdlib.html
			DefaultVarcharSize: 256,                                               // add default size for string fields, not a primary key, no index defined and don't have default values
			DefaultBlobSize:    256,                                               // add default size for bytes fields, not a primary key, no index defined and don't have default values
			DisableDeletion:    true,                                              // disable row deletion, which not supported before ImmuDB 1.2
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		// Migrate the schema
		err = db.AutoMigrate(&Product{})
		if err != nil {
			panic(err)
		}
		// Create
		err = db.Create(&Product{Code: "D43", Price: 100, Amount: 500}).Error
		if err != nil {
			panic(err)
		}
		// Read
		var product Product
		// find just created one
		err = db.First(&product).Error
		if err != nil {
			panic(err)
		}
		// find product with code D42
		err = db.First(&product, "code = ?", "D43").Error
		if err != nil {
			panic(err)
		}
		// Update - update product's price to 200
		err = db.Model(&product).Update("Price", 888).Error
		if err != nil {
			panic(err)
		}

		// Update - update multiple fields
		err = db.Model(&product).Updates(Product{Price: 200, Code: "F42"}).Error
		if err != nil {
			panic(err)
		}

		err = db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"}).Error
		if err != nil {
			panic(err)
		}

		// Delete - delete product
		// err = db.Delete(&product, product.ID).Error
		// if err != nil {
		// 	panic(err)
		// }
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
}
