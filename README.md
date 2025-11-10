# HumanFriendlyPasswordGeneendlyPasswordGenerator

A compiled Go application that generates highly secure, long passphrases designed to be **manually typeable** in high-friction environments (like mobile devices or RDP sessions).

This tool is built on the philosophy that **Length & Usability > Raw Complexity**. A 30-character passphrase you can actually type is more secure than a 20-character random string you can't.

---

## ✨ Features

* **Single, Portable Binary:** Uses `go:embed` to bundle the wordlist directly into the executable. No external dependencies needed.
* **Cross-Platform:** Can be compiled for Windows, macOS, and Linux from any Go environment.
* **Secure by Default:** The default settings generate a passphrase with **~72 bits of entropy**.
* **Quality of Life:**
    * Generate multiple options at once (`-n 5`).
    * Automatically copy the first password to the clipboard (`-c`).
* **Highly Flexible:** Control word count, separator count, digit count (fixed or range), separator symbols, and capitalization rules.
* **Interoperable by Design:** The core wordlist and default symbol set are restricted to characters that are identical on all major keyboard layouts (QWERTZ/QWERTY), preventing RDP/VM typing errors.

---

## 🚀 Usage

The tool is a single command-line executable.

**Generate a password with default settings:**
```bash
./hfpg
````

**Generate 3 passwords, 5 words long, 2 separators, and copy the first one:**

```bash
./hfpg -n 3 -c -w 5 -s 2
```

**See all available options:**

```bash
./hfpg -h
```

-----

## Default Behavior

Running `./hfpg` without flags is equivalent to:

```bash
./hfpg -w 4 -s 1 -sep "!" -d-range 4 -typo-rate 0.33 -caps camel
```

This generates a password with:

  * 4 random words.
  * 1 separator block (`!xxxx!`).
  * The block contains exactly 4 digits.
  * A guaranteed minimum of one "typo" (a swapped character, e.g., `Wrod`).
  * `CamelCase` capitalization.

**Example Default Output:** `WortEins!1234!WortZweiWroDreiWortVier`

-----

## 🔐 Security & Entropy

### How Secure is the Default Password? (\~72 Bits)

We must assume an attacker **knows our method** (Kerckhoffs's Principle). They won't try to brute-force `aaaaa...` (which would be \~179 bits). They will attack the *components*.

The default password (4 words, 1 block) has **\~72 bits of entropy**.

| Component | Calculation (Default) | Entropy (Bits) |
| :--- | :--- | :--- |
| **Words** | 4 words from 13,131: `13131^4` | \~54.7 bits |
| **Digits** | 4 digits: `10^4` | \~13.3 bits |
| **Position** | 1 separator in 3 positions | \~1.6 bits |
| **Typo** | 1 typo in \~5 positions | \~2.3 bits |
| **Total** | | **\~71.9 bits** |

### What is 72 Bits?

  * **vs. Traditional Passwords:** To get 72 bits of entropy with a traditional password (using 95 characters: a-z, A-Z, 0-9, all symbols), you would need **11-12 random characters**.
  * **vs. Cracking (Modern Hashing):** If stored with **bcrypt** or **Argon2**, cracking 72 bits would take a supercomputer cluster **billions of years**.
  * **vs. Cracking (Broken Hashing):** If stored with **MD5** (fast hash), it would still take a supercomputer cluster **months to years** to crack.

### How to Increase Entropy (The Levers)

The parameters have different impacts on security.

1.  **Strongest Levers (Use these):**

      * **`-w` (Word Count):** This is the most powerful lever. Each additional word adds **\~13.7 bits**.
      * **`-s` (Separator Count):** The second most powerful. Each additional separator block adds **\~13.3 bits** (for 4 digits).

2.  **Good Levers:**

      * **`--caps=random`:** Changes from predictable `CamelCase` to `rAndom`. This adds `log2(word_length)` per word, totaling \~10-12 bits for 4 words.

3.  **Weak Levers (Good for rules, not for raw entropy):**

      * **`-d-range`:** Increasing from 4 to 5 digits only adds `log2(10)` (\~3.3 bits).
      * **`-sep`:** Changing from `!` to `+-!_` (4 chars) only adds `log2(4)` (= 2 bits) per separator.

-----

## 🛡️ Core Interoperability Rules

The generator is designed to avoid typing errors in VMs, RDP sessions, or on foreign keyboards.

### 1\. Character Set (Letters)

The default wordlist **excludes** words containing characters that change position:

  * **No `Z` or `Y`**
  * **No Diacritics/Umlauts** (`ä`, `ö`, `ü`, `ß`, etc.)

### 2\. Special Characters (The "Universal Set")

The default separator pool (`+`, `-`, `_`, `!`) only uses characters that are in the same physical location on US (QWERTY) and German (QWERTZ) layouts.

-----

## 💡 Examples

### Default

*Command:* `./hfpg`
*Output:* `WortEins!1234!WortZweiWroDreiWortVier`

### Higher Entropy & QoL

*Command:* `./hfpg -n 2 -c -w 5 -s 2 -sep "+-_!" -d-range 2-4 -caps random`
*Output (copies the first line):*

```
WortEins+12+WoRtZwei_4567_WortDreiWortVierWortFuenf
KafFee!12!AutoBahnReciHte+9876+WortFuenf
```

-----

## 💻 Building from Source

You must have Go 1.16+ (for `go:embed`) installed.

**Build optimized binary:**

```bash
go build -ldflags="-s -w" -o build/hfpg .
```

**Cross-compile for Windows:**

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/hfpg.exe .
```

```
```
rator

A compiled Go application that generates highly secure, long passphrases designed to be **manually typeable** in high-friction environments (like mobile devices or RDP sessions).

This tool is built on the philosophy that **Length & Usability > Raw Complexity**. A 30-character passphrase you can actually type is more secure than a 20-character random string you can't.

---

## ✨ Features

* **Single, Portable Binary:** Uses `go:embed` to bundle the wordlist directly into the executable. No external dependencies needed.
* **Cross-Platform:** Can be compiled for Windows, macOS, and Linux from any Go environment.
* **Secure by Default:** The default settings generate a passphrase with **~72 bits of entropy**.
* **Quality of Life:**
    * Generate multiple options at once (`-n 5`).
    * Automatically copy the first password to the clipboard (`-c`).
* **Highly Flexible:** Control word count, separator count, digit count (fixed or range), separator symbols, and capitalization rules.
* **Interoperable by Design:** The core wordlist and default symbol set are restricted to characters that are identical on all major keyboard layouts (QWERTZ/QWERTY), preventing RDP/VM typing errors.

---

## 🚀 Usage

The tool is a single command-line executable.

**Generate a password with default settings:**
```bash
./hfpg
````

**Generate 3 passwords, 5 words long, 2 separators, and copy the first one:**

```bash
./hfpg -n 3 -c -w 5 -s 2
```

**See all available options:**

```bash
./hfpg -h
```

-----

## Default Behavior

Running `./hfpg` without flags is equivalent to:

```bash
./hfpg -w 4 -s 1 -sep "!" -d-range 4 -typo-rate 0.33 -caps camel
```

This generates a password with:

  * 4 random words.
  * 1 separator block (`!xxxx!`).
  * The block contains exactly 4 digits.
  * A guaranteed minimum of one "typo" (a swapped character, e.g., `Wrod`).
  * `CamelCase` capitalization.

**Example Default Output:** `WortEins!1234!WortZweiWroDreiWortVier`

-----

## 🔐 Security & Entropy

### How Secure is the Default Password? (\~72 Bits)

We must assume an attacker **knows our method** (Kerckhoffs's Principle). They won't try to brute-force `aaaaa...` (which would be \~179 bits). They will attack the *components*.

The default password (4 words, 1 block) has **\~72 bits of entropy**.

| Component | Calculation (Default) | Entropy (Bits) |
| :--- | :--- | :--- |
| **Words** | 4 words from 13,131: `13131^4` | \~54.7 bits |
| **Digits** | 4 digits: `10^4` | \~13.3 bits |
| **Position** | 1 separator in 3 positions | \~1.6 bits |
| **Typo** | 1 typo in \~5 positions | \~2.3 bits |
| **Total** | | **\~71.9 bits** |

### What is 72 Bits?

  * **vs. Traditional Passwords:** To get 72 bits of entropy with a traditional password (using 95 characters: a-z, A-Z, 0-9, all symbols), you would need **11-12 random characters**.
  * **vs. Cracking (Modern Hashing):** If stored with **bcrypt** or **Argon2**, cracking 72 bits would take a supercomputer cluster **billions of years**.
  * **vs. Cracking (Broken Hashing):** If stored with **MD5** (fast hash), it would still take a supercomputer cluster **months to years** to crack.

### How to Increase Entropy (The Levers)

The parameters have different impacts on security.

1.  **Strongest Levers (Use these):**

      * **`-w` (Word Count):** This is the most powerful lever. Each additional word adds **\~13.7 bits**.
      * **`-s` (Separator Count):** The second most powerful. Each additional separator block adds **\~13.3 bits** (for 4 digits).

2.  **Good Levers:**

      * **`--caps=random`:** Changes from predictable `CamelCase` to `rAndom`. This adds `log2(word_length)` per word, totaling \~10-12 bits for 4 words.

3.  **Weak Levers (Good for rules, not for raw entropy):**

      * **`-d-range`:** Increasing from 4 to 5 digits only adds `log2(10)` (\~3.3 bits).
      * **`-sep`:** Changing from `!` to `+-!_` (4 chars) only adds `log2(4)` (= 2 bits) per separator.

-----

## 🛡️ Core Interoperability Rules

The generator is designed to avoid typing errors in VMs, RDP sessions, or on foreign keyboards.

### 1\. Character Set (Letters)

The default wordlist **excludes** words containing characters that change position:

  * **No `Z` or `Y`**
  * **No Diacritics/Umlauts** (`ä`, `ö`, `ü`, `ß`, etc.)

### 2\. Special Characters (The "Universal Set")

The default separator pool (`+`, `-`, `_`, `!`) only uses characters that are in the same physical location on US (QWERTY) and German (QWERTZ) layouts.

-----

## 💡 Examples

### Default

*Command:* `./hfpg`
*Output:* `WortEins!1234!WortZweiWroDreiWortVier`

### Higher Entropy & QoL

*Command:* `./hfpg -n 2 -c -w 5 -s 2 -sep "+-_!" -d-range 2-4 -caps random`
*Output (copies the first line):*

```
WortEins+12+WoRtZwei_4567_WortDreiWortVierWortFuenf
KafFee!12!AutoBahnReciHte+9876+WortFuenf
```

-----

## 💻 Building from Source

You must have Go 1.16+ (for `go:embed`) installed.

**Build optimized binary:**

```bash
go build -ldflags="-s -w" -o build/hfpg .
```

**Cross-compile for Windows:**

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/hfpg.exe .
```

```
```
