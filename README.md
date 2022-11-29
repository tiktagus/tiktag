# TikTag

`TikTag` is a command-line tool for preparing images for blog post or sharing.

`TikTag` is shipped as a local command-line app written in Go.

## Features (v0.1)

Tiktag server offers the following features,

1. Install and run as a local command line app
2. Upload a photo and get its S3 URL back as a response, for use in Markdown for publishing
3. Configurable params via `config.yml`, on  
   * which S3 compatible object storage service to host the photo, default is MinIO
   * configuring file extention, i.e., `png`, `jpg`, etc.
   * initializing user config in `immudb`, for new installation
   * how search is configured, like searching by `filename` or `ttid`

### Commands (the verbs)

1. Host a photo or file,
   
   ```
   > tiktag myfilename.png
   > Success! Here is your hosted asset's URL,
   > https://s3.tikoly.com/village/563583552944996352.png
   ```

2. Search for a stored file and retrieve it's URL,
   
   ```
   > tiktag query myfilename.png
   > Success! Here is your hosted asset's URL,
   > https://s3.tikoly.com/village/563583552944996352.png
   ```
   
   if your asset not found,
   
   ```
   > Oops...asset not found, or run `tiktag <myfile>` to host it
   ```

## Design

Tiktag was firstly designed to streamline photo preparation before being hosted and published/referenced in a blog post. We'll see where it leads us.

### Data objects (key nouns)

List of key data objects in TikTag,

* `ttasset`, top level noun / object with properties below,
  * `ttid`, TikTag ID, unique ID for each file/image/object, also as its filename, i.e., `563583552944996352.png`
  * `ttidhash`, the hash generated upon initial upload of a file, for cryptographic verification of the file's integrity
  * `filename`, original filename of a photo/file, whatever it is
  * `fileext`, file extension name, such as `.png`, `.jpg`
  * `tturl`, TikTag URL of a hosted file/photo, i.e., `https://s3.tikoly.com/village/563583552944996352.png`
    * it's constructed like, `{TargetURL}/{TargetBucket}/{ttid}.{fileext}`
* S3 (MinIO) related,
  * `TargetURL`
  * `AccessKey`
  * `SecretKey`
  * `TargetBucket`

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

## Tech Stack (to evaluate)

Candidate tech dependencies for making TikTag happen,

* Written in Go
* Data store, [ImmuDB](https://github.com/codenotary/immudb), enforcing immutable data policies
  * more about `ImmuDB` on [this podcast](https://changelog.com/gotime/219).
* Object storage, S3 API compatible object storage (hosted on [Minio](https://github.com/minio/minio) by default)
  * RClone, for handling jobs with object storage
* (maybe) Background (cron) job handling on [Dkron](https://dkron.io/), and [GoFS](https://github.com/no-src/gofs), across-platform file synchronization tool out of the box based on golang

## Contributors

* [Chance Jiang](https://github.com/chancefcc), designer
* [Atman An](https://github.com/twinsant), lead developer
* Ryan Sy, lead on user success 

## Licensing

[The 3-Clause BSD License](https://opensource.org/licenses/BSD-3-Clause)
