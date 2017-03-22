---
layout: "docs"
page_title: "Using PGP, GPG, and Keybase"
sidebar_current: "docs-concepts-pgp-gpg-keybase"
description: |-
  Vault has the ability to integrate with OpenPGP-compatible programs like GPG
  and services like Keybase.io to provide an additional layer of security when
  performing certain operations.  This page details the various GPG
  integrations, their use, and operation.
---

# Using PGP, GPG, and Keybase

Vault has the ability to integrate with OpenPGP-compatible programs like GPG
and services like Keybase.io to provide an additional layer of security when
performing certain operations.  This page details the various PGP integrations,
their use, and operation.

## Initializing with PGP
One of the early fundamental problems when bootstrapping and initializing Vault
was that the first user (the initializer) received a plain-text copy of all of
the unseal keys. This defeats the promises of Vault's security model, and it
also makes the distribution of those keys more difficult. Since Vault 0.3,
Vault can optionally be initialized using PGP keys. In this mode, Vault will
generate the unseal keys and then immediately encrypt them using the given
users' public PGP keys. Only the owner of the corresponding private key is then
able to decrypt the value, revealing the plain-text unseal key.

First, you must create, acquire, or import the appropriate key(s) onto the
local machine from which you are initializing Vault. This guide will not
attempt to cover all aspects of PGP keys but give examples using two popular
programs: Keybase and GPG.

For beginners, we suggest using [Keybase.io](https://keybase.io/) ("Keybase")
as it can be both simpler and has a number of useful behaviors and properties
around key management, such as verification of users' identities using a number
of public online sources. It also exposes the ability for users to have PGP
keys generated, stored, and managed securely on their servers. Using Vault with
Keybase will be discussed first as it is simpler.

## Initializing with Keybase
To generate unseal keys for Keybase users, Vault accepts the `keybase:` prefix
to the `-pgp-keys` argument:

```
$ vault init -key-shares=3 -key-threshold=2 \
    -pgp-keys="keybase:jefferai,keybase:vishalnayak,keybase:sethvargo"
```

This requires far fewer steps that traditional PGP (e.g. with `gpg`) because
Keybase handles a few of the tedious steps. The output will be the similar to
the following:

```
Key 1: c1c04c03d5f43b6432ea77f3010800...
Key 2: 612b611295f255baa2eb702a5e254f...
Key 3: ebfd78302325e2631bcc21e11cae00...
...
```

### Unsealing with Keybase
As a user, the easiest way to decrypt your unseal key is with the Keybase CLI
tool. You can download it from [Keybase.io download
page](https://keybase.io/download). After you have downloaded and configured
the Keybase CLI, you are now tasked with entering your unseal key. To get the
plain-text unseal key, you must decrypt the value given to you by the
initializer. To get the plain-text value, run the following command:

```
$ echo "c1c0..." | xxd -r -p | keybase pgp decrypt
```

And replace `c1c0...` with the encrypted key.

You will be prompted to enter your Keybase passphrase. The output will be the
plain-text unseal key.

```
6ecb46277133e04b29bd0b1b05e60722dab7cdc684a0d3ee2de50ce4c38a357101
```

This is your unseal key in plain-text and should be guarded the same way you
guard a password. Now you can enter your key to the `unseal` command:

```
$ vault unseal
Key (will be hidden): ...
```

- - -

## Initializing with GPG
GPG is an open-source implementation of the OpenPGP standard and is available
on nearly every platform. For more information, please see the [GPG
manual](https://gnupg.org/gph/en/manual.html).

To create a new PGP key, run, following the prompts:

```
$ gpg --gen-key
```

To import an existing key, download the public key onto disk and run:

```
$ gpg --import key.asc
```

Once you have imported the users' public keys, you need to save their values
to disk as either base64 or binary key files. For example:

```
$ gpg --export 348FFC4C | base64 > seth.asc
```

These key files must exist on disk in base64 or binary. Once saved to disk, the
path to these files can be specified as an argument to the `-pgp-keys` flag.

```
$ vault init -key-shares=3 -key-threshold=2 \
    -pgp-keys="jeff.asc,vishal.asc,seth.asc"
```

The result should look something like this:

```
Key 1: c1c04c03d5f43b6432ea77f3010800...
Key 2: 612b611295f255baa2eb702a5e254f...
Key 3: ebfd78302325e2631bcc21e11cae00...
...
```

The output should be rather long in comparison to a regular unseal key. These
keys are encrypted, and only the user holding the corresponding private key
can decrypt the value. The keys are encrypted in the order in which specified
in the `-pgp-keys` attribute. As such, the first key belongs to Jeff, the second
to Vishal, and the third to Seth. These keys can be distributed over almost any
medium, although common sense and judgement are best advised.

### Unsealing with a GPG
Assuming you have been given an unseal key that was encrypted using your public
PGP key, you are now tasked with entering your unseal key. To get the
plain-text unseal key, you must decrypt the value given to you by the
initializer. To get the plain-text value, run the following command:

```
$ echo "c1c0..." | xxd -r -p | gpg -d
```

And replace `c1c0...` with the encrypted key.

If you encrypted your private PGP key with a passphrase, you may be prompted to
enter it.  After you enter your password, the output will be the plain-text
key:

```
6ecb46277133e04b29bd0b1b05e60722dab7cdc684a0d3ee2de50ce4c38a357101
```

This is your unseal key in plain-text and should be guarded the same way you
guard a password. Now you can enter your key to the `unseal` command:

```
$ vault unseal
Key (will be hidden): ...
```
