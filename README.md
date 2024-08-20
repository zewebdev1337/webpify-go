# webpify-go

Go program to convert all PNG images inside a folder to WebP format, multithreaded for maximum CPU utilization and speed.

## Features

- **Multithreaded:** Utilizes all available CPU cores.
- **Customizable Settings:** Adjust output scale and WebP quality.
- **Removes Transparency Data:** Removes alpha channel information for saving a few extra KB.

## Usage

1. Download binary from releases and place somewhere recognized by PATH.

2. Call from within the folder with the PNGs.

    $ `webpify`

3. Enter the desired output scale (e.g., 0.5 for 50% reduction) when prompted.

4. Enter the desired WebP quality (e.g., 80) when prompted.

Converted files will have the same name as the original PNG files with a `.webp` extension.