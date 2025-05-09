
# shellfinder

A Simple Tool to Find Shells and Some Interesting Endpoints in Websites

A shell is malicious PHP file and execute it by accessing it via a web browser. It is a PHP script that allows the attacker to control the server - essentially a backdoor program, similar in functionality to a trojan for personal computers.

## Installation 
Prerequisites:

- Python 3.6

- requests module

1. Put your websites in `websites.txt`

2. Install requests `pip3 install requests`

3. Run the script `$ python3 shellfinder.py`

## Contribution
Please feel free to fork the repository and make pull requests.

Made with love as always.


# Go Shell Finder

A fast, concurrent web shell finder written in Go. This tool scans a list of websites against a list of common (or custom) web shell endpoints to identify potential compromises by checking for `200 OK` HTTP responses.

Based on the original concept of "Shell Finder v0.0.1 by Ahmed Lekssays (0x70776e)" (Python version). This Go version aims for significantly improved performance through concurrency.

## Features

*   **Concurrent Scanning:** Utilizes Go routines to scan multiple URLs simultaneously, making it much faster than sequential scanners.
*   **Input Files:** Reads target websites from `websites.txt` and shell endpoints from `endpoints.txt`.
*   **Customizable Concurrency:** The number of concurrent workers can be easily adjusted in the source code (`numWorkers` constant).
*   **HTTP Timeout:** Configurable HTTP request timeout to prevent indefinite hangs (`Timeout` in `http.Client`).
*   **Clear Output:** Prints success messages indicating the full URL where a `200 OK` was received and the specific endpoint that was hit.
*   **User-Agent:** Sets a custom User-Agent for requests.
*   **Basic Error Handling:** Handles common errors like file not found and network issues gracefully for individual requests.

## Prerequisites

*   Go (version 1.16 or newer recommended) installed on your system.

## Installation & Setup

1.  **Clone the repository (or download the `shellfinder.go` file):**
    ```bash
    git clone https://github.com/yourusername/go-shell-finder.git
    cd go-shell-finder
    ```
    (Replace `yourusername/go-shell-finder` with your actual repository path if you create one.)
    Alternatively, just download or copy the `shellfinder.go` file into a directory.

2.  **Prepare Input Files:**
    Create two text files in the same directory as `shellfinder.go`:

    *   **`websites.txt`**:
        Each line should contain a full website URL (including `http://` or `https://`).
        Do **not** end URLs with a `/`.
        Example `websites.txt`:
        ```
        http://example.com
        https://test-site.org
        http://another-target.net
        ```

    *   **`endpoints.txt`**:
        Each line should contain a potential shell path/endpoint.
        Leading slashes `/` will be trimmed if present.
        Example `endpoints.txt`:
        ```
        shell.php
        c99.php
        r57.php
        uploads/shell.php
        wp-content/uploads/shell.php
        admin/shell.php
        webshell.aspx
        cmd.jsp
        ```

## Method of Use

1.  **Navigate to the directory** containing `shellfinder.go`, `websites.txt`, and `endpoints.txt`.

2.  **Build the executable (recommended):**
    ```bash
    go build shellfinder.go
    ```
    This will create an executable file named `shellfinder` (or `shellfinder.exe` on Windows).

3.  **Run the scanner:**
    ```bash
    ./shellfinder
    ```
    On Windows:
    ```bash
    .\shellfinder.exe
    ```

    Alternatively, you can run directly without building (slower startup):
    ```bash
    go run shellfinder.go
    ```

## Example Output

The script will first print information about loading websites and endpoints. Then, as it finds potential shells (URLs returning a `200 OK` status), it will print:
