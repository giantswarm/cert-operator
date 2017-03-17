# Usage
Here we describe how to use `certctl` for various operations. For Vault
specific documentation please have a look at https://www.vaultproject.io/docs.

At first you want to prepare your environment to provide the Vault address and
the Vault token to `certctl`. You can also provide these settings using the
command line flags, but then your commands are getting more messed up. Note
that you need to set a token having root access for the upcoming cluster setup.
```
export VAULT_ADDR=<vault-addr>
export VAULT_TOKEN=<vault-root-token>
```

When you want to know the state of a cluster, use the `inspect` command. Here
we see there had no setup happen yet.
```
$ certctl inspect --cluster-id=123
Inspecting cluster for ID '123':

    PKI backend mounted: false
    Root CA generated:   false
    PKI role created:    false
    PKI policy created:  false

Tokens may have been generated for this cluster. Created tokens
cannot be shown as they are secret. Information about these
secrets needs to be looked up directly from the location of the
cluster's installation.
```

Setting up a cluster works using the `setup` command. It is shown what happend.
`setup` can be called multiple times. A PKI backend is only mounted if it is
not mounted yet. A root CA is only generated if it is not generated yet. You
get the picture. One exception of this behaviour is the token generation. Each
call to `setup` generates `--num-tokens` tokens. So in case you need one more
token it is safe to simply call `setup` again for a specific cluster ID. Note
that we set a Vault token with root capabilities for the cluster setup.
```
$ certctl setup --allowed-domains=giantswarm.io --common-name=giantswarm.io --cluster-id=123
Set up cluster for ID '123':

    - PKI backend mounted
    - Root CA generated
    - PKI role created
    - PKI policy created

The following tokens have been generated for this cluster:

    265c159b-0fee-4198-b6d5-b6416752be04

```

When we now call `inspect` again we see that the cluster is set up properly.
```
$ certctl inspect --cluster-id=123
Inspecting cluster for ID '123':

    PKI backend mounted: true
    Root CA generated:   true
    PKI role created:    true
    PKI policy created:  true

Tokens may have been generated for this cluster. Created tokens
cannot be shown as they are secret. Information about these
secrets needs to be looked up directly from the location of the
cluster's installation.
```

In case the cluster is set up, we can generate certificates for it using the
`issue` command. Note that `issue` should only be provided the restricted token
generated on `setup`. That way it is more safe to automate the certificate
generation. Any process can then issue certificates for the cluster it was set
up for, and not more.
```
export VAULT_TOKEN=<vault-cluster-token>
```

```
certctl issue --cluster-id=123 --common-name=admin.giantswarm.io --crt-file=./crt.pem --key-file=./key.pem --ca-file=./ca.pem
Public key written to './crt.pem'.
Private key written to './key.pem'.
Root CA written to './ca.pem'.
```

At some point a cluster may not be used anymore, or needs to be cleaned up for
some reason. Here we can use the `cleanup` command. Note that a root token is
again necessary to cleanup a cluster.
```
export VAULT_TOKEN=<vault-root-token>
```

```
$ certctl cleanup --cluster-id=123
Inspecting cluster for ID '123':

    - PKI backend unmounted
    - Root CA deleted
    - PKI role deleted
    - PKI policy deleted

Tokens may have been generated for this cluster. Created tokens
cannot be revoked here as they are secret. Tokens need to be
revoked manually. In case a cluster with the same ID will be
generated, tokens generated for this cluster will be able to
access this new cluster again. Information about these secrets
needs to be looked up directly from the location of the cluster's
installation.
```

When we now inspect the cluster again, we see that it is no longer set up.
```
$ certctl inspect --cluster-id=123
Inspecting cluster for ID '123':

    PKI backend mounted: false
    Root CA generated:   false
    PKI role created:    false
    PKI policy created:  false

Tokens may have been generated for this cluster. Created tokens
cannot be shown as they are secret. Information about these
secrets needs to be looked up directly from the location of the
cluster's installation.
```
