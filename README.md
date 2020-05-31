# GSC Utils

## Usage

### Uploader

```bash
$ gcs-upload --help
$ gcs-upload --bucket bucket --input input.csv --object object
```

### Downloader

```bash
$ gcs-download --help
$ gcs-download --bucket bucket --object object --output output.csv
```

## Installation

### Linux

```bash
$ curl -Ls https://github.com/pottava/gcs-utils/releases/download/upload-v0.5.0/gcs-upload-linux \
  -o gcs-upload && chmod +x gcs-upload
$ curl -Ls https://github.com/pottava/gcs-utils/releases/download/download-v0.5.0/gcs-download-linux \
  -o gcs-download && chmod +x gcs-download
```

### Max

```bash
$ curl -Ls https://github.com/pottava/gcs-utils/releases/download/upload-v0.5.0/gcs-upload-mac \
  -o gcs-upload && chmod +x gcs-upload
$ curl -Ls https://github.com/pottava/gcs-utils/releases/download/download-v0.5.0/gcs-download-mac \
  -o gcs-download && chmod +x gcs-download
```

### Windows

```bash
$ 
```
