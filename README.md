# UTime - Go Time Utility

A command-line utility written in Go for displaying various time information, including current local/UTC time, Unix epoch time, Swatch Internet Time (@beats), and converting between epoch/beat time and local time.

## Features

*   **Current Time Display:** Shows the current time in both UTC and the system's local timezone, including timezone information.
*   **Epoch Time:** Displays the current Unix epoch time (seconds since Jan 1, 1970 UTC).
*   **Swatch Internet Time (@beats):** Displays the current time in Swatch Internet Time.
*   **Epoch Conversion:** Converts a given Unix epoch timestamp into the system's local date and time.
*   **Beat Time Conversion:** Converts a given Swatch Internet Time (@beat) into the corresponding local time range for the current day.

## Installation

You need to have Go installed on your system.

1.  **Clone the repository (or download the source code):**
    ```bash
    # If you have a git repo (replace with your actual repo URL)
    git clone https://github.com/yourusername/utime-go.git
    cd utime-go
    ```
    Or simply navigate to the directory containing `utime.go`.

2.  **Build the executable:**
    ```bash
    go build -o utime utime.go
    ```
    This will create an executable file named `utime` (or `utime.exe` on Windows) in the current directory.

3.  **(Optional) Move to your PATH:**
    You can move the `utime` executable to a directory in your system's PATH (like `/usr/local/bin` or `~/bin`) to run it from anywhere.
    ```bash
    # Example for Linux/macOS
    sudo mv utime /usr/local/bin/
    ```

## Usage

Usage: utime [-epoch] [-beat] [@<beat_time> | <epoch_time>]
Displays current UTC and local time information by default.
Options:
-epoch Print the current Unix epoch time.
-beat Print the current Swatch Internet Time (@beats).
<epoch_time> Convert the given Unix epoch time (integer seconds) to local time.
@<beat_time> Convert the given Swatch Internet Time (0-999) to a local time range for the current day.
-h, --help Show this help message.


*   Running `utime` with no arguments displays the default current time information.
*   Flags (`-epoch`, `-beat`) show specific current time formats. These flags cannot be combined with each other or with time arguments.
*   Providing a numeric argument interprets it as a Unix epoch time to be converted to local time.
*   Providing an argument starting with `@` followed by a number (0-999) interprets it as a Swatch Internet Time beat to be converted into a local time range.

## Examples

1.  **Show current UTC and Local Time:**
    ```bash
    ./utime
    ```
    *Output (example):*
    ```
    Current Time Information:
    ==============================
    UTC Time      : 2023-10-27 10:30:00 UTC
    Local Time    : 2023-10-27 06:30:00 EDT
    Time Zone     : EDT (UTC-4)
    Location      : America/New_York
    ==============================
    ```

2.  **Show current Unix Epoch Time:**
    ```bash
    ./utime -epoch
    ```
    *Output (example):*
    ```
    1698402600
    ```

3.  **Show current Swatch Internet Time (@beats):**
    ```bash
    ./utime -beat
    ```
    *Output (example):*
    ```
    @479.16
    ```

4.  **Convert a specific Epoch time to Local Time:**
    (Epoch 1678886400 corresponds to 2023-03-15 12:00:00 UTC)
    ```bash
    ./utime 1678886400
    ```
    *Output (example, assuming EDT timezone):*
    ```
    2023-03-15 08:00:00 EDT
    ```

5.  **Convert a specific Beat time to a Local Time range:**
    ```bash
    ./utime @500
    ```
    *Output (example, assuming EDT timezone on 2023-10-27):*
    ```
    @500 corresponds to the time range 08:00:00.000 - 08:01:26.400 EDT (on 2023-10-27)
    ```

6.  **Show Help:**
    ```bash
    ./utime -h
    # or
    ./utime --help
    ```
    *Output:* Displays the usage information.

## Contributing

Feel free to open issues or submit pull requests if you find bugs or have suggestions for improvements.
