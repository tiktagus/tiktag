# TikTag

`TikTag` is a command-line tool for hosting and sharing digital asssets via object storage.

## Features (v0.1)

Tiktag server offers the following features,

1. Install and run as a local command line app
2. Upload and host a photo to your object storage of choice, and get its S3 URI, for publishing or sharing
3. Configurable params via `config.yml`, such as,
   * your S3 compatible object storage service of choice, such as your own MinIO instance;
   * how search is configured;
   * how much asset handling logs you want to share with others

### Install `tiktag`

Before you install `tiktag`,
1. you've got Golang installed on your local environment, typically the Terminal app in MacOS;
2. you've `cd` to your asset's directory, where you're going to use `tiktag` for the job;

```
> go install https://github.com/tikoly-com/tiktag@latest
```

After you've installed `tiktag`,
3. make sure you copy `[config.yaml.sample](https://github.com/tikoly-com/tiktag/blob/main/config.yaml.sample)` to `config.yaml` your asset's directory, and configure your S3-compatible object storage of choice, such as MinIO,

```
 minio:
  endpoint: "s3.example.com"
  accessKey: "example"
  secretKey: "example"
  useSSL: true
  bucketName: "example"
```

### Command (examples)

1. Host a photo or file,
   
   ```
   > tiktag myfilename.png
   > Success! Here is your hosted asset's URL,
   > https://s3.tikoly.com/village/563583552944996352.png
   ```

2. Search for a stored file and retrieve it's URL,
   
   ```
   > tiktag search myfilename.png
   > Success! Here is your hosted asset's URL,
   > https://s3.tikoly.com/village/563583552944996352.png
   ```
   
   if your asset not found,
   
   ```
   > Oops...asset not found, or run `tiktag <myfile>` to host it
   ```

## Design

Tiktag was firstly designed to streamline photo preparation before being hosted and published/referenced in a blog post. We're increasingly seeing Tiktag's potential in NFT-related business applications we're building for our clients world-wide.

We're excited to see where it leads us.

### Data objects (key nouns)

List of key data objects in TikTag,

* `ttasset`, top level noun / object with properties below,
  * `ttid`, TikTag ID, unique ID for each file/image, also as its filename, i.e., `563583552944996352.png`
  * `ttidhash`, the hash generated upon initial upload of a file, for cryptographic verification of the file's integrity
  * `filename`, original filename of a photo/file, whatever it is
  * `fileext`, file extension name, such as `.png`, `.jpg`
  * `tturl`, TikTag URL of a hosted file/photo, i.e., `https://s3.tikoly.com/village/563583552944996352.png`
    * it's constructed like, `{TargetURL}/{TargetBucket}/{ttid}.{fileext}`
* S3 (MinIO) related,
  * `endpoint`
  * `accessKey`
  * `secretKey`
  * `useSSL`
  * `bucketName`

### Key verbs

* `tiktag`, command for tagging and storing an asset
  * example of tagging an asset, `tiktag myfilename.png -b s3aws`
  * example of minting an asset, `tiktag myfilename.png -tz sui` 
    - `tz` is initial of Chance's former/deceased co-founder, Tao Zui, in memory of his ingenuity inspiring Chance's design aesthetics)

## Roadmap

1. v0.1, Validating the idea and usability, a local command line tool
  * jobs are handled locally by user
  * supports major `s3` compatible storage
  * accepts Gtihub Sponsor
2. v1.x, Hosted version,  web GUI and Web3 login (Sui Wallet / testnet, devnet)
  * for tracking assets handling
  * for more non-technical audience
  * better performance

## Tech Stack

Candidate tech dependencies for making TikTag happen,

* Written in Go
* Local data store, [ImmuDB](https://github.com/codenotary/immudb), enforcing immutable data policies
  * more about `ImmuDB` on [this podcast](https://changelog.com/gotime/219).
* Object storage, S3-compatible object storage ([MinIO](https://github.com/minio/minio) by default)

## Contributors

* [Chance Jiang](https://github.com/chancefcc), designer
* [Atman An](https://github.com/twinsant), lead developer
* Ryan Sy, lead on user success 

## Licensing

[The 3-Clause BSD License](https://opensource.org/licenses/BSD-3-Clause)
