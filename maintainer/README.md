# Relay Maintainer

This simple python app maintains the header relay by querying a bcoin node and
pushing headers to an associated geth or infura node.

Generally it follows the crash-only programming paradigm. Rather than
recovering from errors, we expose them, crash, and emphasize safe resume via a
reboot.

It is ALPHA-quality software at best. It is not long-term stable. It has poor
handling of bitcoin reorgs, for example.

## Setup

install `pipenv` and `pyenv`

```sh
$ pipenv install --python=$(pyenv which python3.7)
```

## Testing
Current testscript runs linting and typechecking only.

```sh
$ pipenv run test
```

## Running the header forwarder

Make a config `.env` file in `maintainer/config`.

```sh
$ cp maintainer/config/.sample.env maintainer/config/.my_env_file.env

# update the env to point to your BCOIN node, and either geth or infura
$ vim maintainer/config/.env

$ pipenv run python maintainer/header_forwarder/h.py .my_env_file.env
```
