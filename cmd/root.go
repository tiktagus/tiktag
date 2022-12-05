/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

func fakeMachineID(uint16) bool {
	return true
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
		ctx := context.Background()
		endpoint := "s3.tikoly.com"
		// FIXME: DO NOT commit!!!!
		accessKeyID := "ZaRTBCf2g4ZMGVgu"
		secretAccessKey := "xCHRw9vv1etAUx0pABvRQiDLnF4euowj"
		useSSL := true

		// Initialize minio client object.
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}

		// Make a new bucket called mymusic.
		bucketName := "tiktag"

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
		objectName := "demo.jpg"
		filePath := "./demo.jpg"
		contentType := "image/jpeg"

		// Upload the zip file with FPutObject
		info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
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
