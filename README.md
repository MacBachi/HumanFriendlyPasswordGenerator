# HumanFriendlyPasswordGenerator

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

| Component    | Calculation (Default)          | Entropy (Bits)  |
| :----------- | :----------------------------- | :-------------- |
| **Words**    | 4 words from 13,131: `13131^4` | \~54.7 bits     |
| **Digits**   | 4 digits: `10^4`               | \~13.3 bits     |
| **Position** | 1 separator in 3 positions     | \~1.6 bits      |
| **Typo**     | 1 typo in \~5 positions        | \~2.3 bits      |
| **Total**    |                                | **\~71.9 bits** |

### What is 72 Bits?

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

## 📊 Security Context: How This Compares

### vs. Traditional Random Passwords (e.g., `k!8(Gf%z#p@L`)

A traditional 11-12 character random password (using 95 characters) has the **same \~72 bits of entropy**.

  * **The Problem:** `k!8(Gf%z#p@L` is **impossible to type** from memory. You cannot enter it on a mobile device or RDP session without extreme effort and high error rates.
  * **The Failure Mode:** Because it's unusable, users write it on a sticky note (insecure) or revert to a weaker password they *can* remember.
  * **The Solution:** `WortEins!1234!WortZweiWroDreiWortVier` provides the **exact same mathematical security** as `k!8(Gf%z#p@L`, but is delivered in a "chunkable" format that a human can actually read, verify, and type.

### vs. Common "Real-World" Passwords (e.g., `Password123!`)

Most data breaches are not caused by brute-forcing 72-bit passwords. They are caused by users picking simple, predictable patterns.

  * **The Problem:** Users choose passwords like `Summer2024!`, `Mustang69`, `Qwertz123`, or `P@ssw0rd!`.
  * **The Attack:** Attackers do not brute-force these. They use **dictionary and mask attacks**. They take massive lists of common words (like `RockYou.txt`) and apply common mutations (like adding `123` or `!` at the end).
  * **The Result:** These passwords have extremely low entropy (20-30 bits) and are often cracked in **seconds**.
  * **The Solution:** This generator defeats dictionary attacks by using a **long combination of *truly random* words**, augmented with random numbers and typos. An attacker's dictionary (which contains `Password`) does not contain the random combination `HausBootSonneAuto`.

-----

## 🛡️ Final Security Recommendations

A strong password is one layer of defense, but it is not enough on its own.

### 1\. Use Multi-Factor Authentication (MFA / 2FA)

MFA is the single most important security measure you can take. It means that even if an attacker *steals* your password (whether it's `123456` or the 72-bit one from this generator), they **cannot log in** because they do not have your second factor (e.g., your phone app, USB key).

**Always enable MFA on every service that supports it.**

### 2\. Use a Unique Password for Every Service

Never reuse passwords. Attackers use a technique called **Credential Stuffing**.

1.  A website ("BadHats.com") gets breached.
2.  Attackers download the user list (e.g., `user@email.com` : `Wort!1234!WortZwei`).
3.  They take that *exact* email and password combination and try it on your bank, your email provider, and your social media.
4.  If you reused your password, they are now in.

**This tool (or a password manager) should be used to create a *different*, strong password for every single account you own.**

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