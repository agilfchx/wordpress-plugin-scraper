# PoC WordPress Vulnerability Scanner

![GitHub release](https://img.shields.io/github/v/release/kavalkala/poc-wordpress-vulnerability-scanner)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/kavalkala/poc-wordpress-vulnerability-scanner/ci.yml?branch=main)
![GitHub issues](https://img.shields.io/github/issues/kavalkala/poc-wordpress-vulnerability-scanner)
![GitHub stars](https://img.shields.io/github/stars/kavalkala/poc-wordpress-vulnerability-scanner)
![GitHub forks](https://img.shields.io/github/forks/kavalkala/poc-wordpress-vulnerability-scanner)
![GitHub license](https://img.shields.io/github/license/kavalkala/poc-wordpress-vulnerability-scanner)
![Twitter](https://img.shields.io/twitter/follow/kavalkala?style=social)

## Description

This tool scans WordPress plugins for PHP CodeSniffer (PHPCS) errors and warnings according to WordPress Coding Standards. It outputs PHP files with detected errors or warnings for further analysis, excluding syntax-related exceptions.

## Features

- Scans for PHP CodeSniffer errors and warnings
- WordPress Coding Standards compliant
- Outputs detected issues for further analysis
- Excludes syntax-related exceptions

## Usage

1. **Clone this repository**:
    ```sh
    git clone https://github.com/your-username/poc-wordpress-vulnerability-scanner.git
    cd poc-wordpress-vulnerability-scanner
    ```

2. **Install dependencies**:
    ```sh
    composer install
    ```

3. **Run the scanner**:
    ```sh
    php scan.php
    ```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## Code of Conduct

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## License

This project is licensed under the terms of the MIT license. See the [LICENSE](LICENSE) file for details.
