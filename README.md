# Alan

[![License Apache 2][badge-license]](LICENSE)
[![GitHub version](https://badge.fury.io/gh/nlamirault%2Falan.svg)](https://badge.fury.io/gh/nlamirault%2Falan)

* Master : [![pipeline status](https://gitlab.com/nicolas-lamirault/alan/badges/master/pipeline.svg)](https://gitlab.com/nicolas-lamirault/alan/commits/master)
* Develop : [![pipeline status](https://gitlab.com/nicolas-lamirault/alan/badges/develop/pipeline.svg)](https://gitlab.com/nicolas-lamirault/alan/commits/develop)

Alan is a bridge between [Hashicorp Vault](https://www.vaultproject.io/) and some password managers :

* [ ] KeepassXC
* [ ] 1password.com
* [ ] Lastpass
* [ ] Pwsafe

## Installation

You can download the binaries :

* Architecture i386 [ [linux](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_linux_386) / [darwin](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_darwin_386) / [freebsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_freebsd_386) / [netbsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_netbsd_386) / [openbsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_openbsd_386) / [windows](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_windows_386.exe) ]
* Architecture amd64 [ [linux](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_linux_amd64) / [darwin](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_darwin_amd64) / [freebsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_freebsd_amd64) / [netbsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_netbsd_amd64) / [openbsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_openbsd_amd64) / [windows](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_windows_amd64.exe) ]
* Architecture arm [ [linux](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_linux_arm) / [freebsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_freebsd_arm) / [netbsd](https://bintray.com/artifact/download/nlamirault/oss/alan-0.1.0_netbsd_arm) ]


## Usage

* CLI help:

        $ alan help

### Local usage

* Start a Vault development server :

        $ vault server -dev

* Setup Vault :

        $ export VAULT_ADDR="http://localhost:8200"
        $ vault auth list
        Path      Type     Description
        ----      ----     -----------
        token/    token    token based credentials

        $ vault auth enable userpass
        Success! Enabled userpass auth method at: userpass/

        $ vault auth list
        Path         Type        Description
        ----         ----        -----------
        token/       token       token based credentials
        userpass/    userpass    n/a

        $ vault secrets list
        Path          Type         Description
        ----          ----         -----------
        cubbyhole/    cubbyhole    per-token private secret storage
        identity/     identity     identity store
        secret/       kv           key/value secret storage
        sys/          system       system endpoints used for control, policy and debugging

        $ vault policy write alan-policy -<<EOF
        path "secret/*" {
                capabilities = ["create", "read", "update", "delete", "list"]
        }
        EOF

        $ vault policy list
        alan-policy
        default
        root

        $ vault write auth/userpass/users/alan password=turing policies=alan-policy
        Success! Data written to: auth/userpass/users/alan

        $ vault login -method=userpass username=alan password=turing
        Success! You are now authenticated. The token information displayed below
        is already stored in the token helper. You do NOT need to run "vault login"
        again. Future Vault requests will automatically use this token.

        Key                    Value
        ---                    -----
        token                  15589767-1e25-6c44-e8c2-9b6c3ac13099
        token_accessor         5fefe9fe-6da7-b67b-a8f0-47583488057e
        token_duration         768h
        token_renewable        true
        token_policies         [alan-policy default]
        token_meta_username    foo

        $ vault write secret/foo value=yes
        Success! Data written to: secret/foo
        $ vault read secret/foo
        Key                 Value
        ---                 -----
        refresh_interval    768h
        value               yes

* Display database entries :

        $ alan keepassxc show --database alan.kdbx
        Please input your password:
        Dev
        Github: foo https://github.com
        Gitlab: foo https://gitlab.com
        Social
        Twitter: alan https://twitter.com
        >>> foo https://fake.social
        Root

* Import a KeepassXC database into the Vault:

        $ alan keepassxc import --database alan.kdbx
        Please input your password:
        Add secret: Dev/Github
        Add secret: Dev/Gitlab
        Add secret: Social/Twitter

* Check entries :

        $ alan vault list
        - Dev/
        - Social/

        $ alan vault list --path Dev
        - Github
        - Gitlab

* Retrieve a secret :

        $ alan vault get --path Dev/Github
        Username: foo
        Password: bar
        URL: https://github.com


## Development

* Initialize environment

        $ make init

* Build tool :

        $ make build

* Launch unit tests :

        $ make test

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md).


## License

See [LICENSE](LICENSE) for the complete license.


## Changelog

A [changelog](ChangeLog.md) is available


## Contact

Nicolas Lamirault <nicolas.lamirault@gmail.com>

[badge-license]: https://img.shields.io/badge/license-Apache2-green.svg?style=flat
