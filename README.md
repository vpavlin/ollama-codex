# ollama-codex

Minimalist CLI tool to use Codex as a store for Ollama models. It allows you to push to and pull from a (local) Codex node. 

## How to use

The simplest workflow to test this tool is to use `push` and then `pull` from the local node (this obviously does not use Codex network and only serves as a minimal functional test)

```
make collama
```

After build succeeds you can find the binary in `./build/collama`

You will first ned to run Codex - follow the official guides or run the `docker-compose.yml`

```
PRIV_KEY=$(openssl rand --hex 32) docker compose up -d
```

Make sure `ollama` is installed and you have a model pulled - e.g. `qwen:0.5b`

```
ollama pull qwen:0.5b
```

Depending on your configuration and Ollama version, your models will be stored either in `~/.ollama/models` or `/usr/share/ollama/modles`. You should be able to use the binary (`./build/collama`) directly if the models are stored in your `HOME`. If they are in `/usr/share` we need to switch to the `ollama` user. Following commands assume the models are stored in `/usr/share`

Then run the `collama.sh push`

```
./collama.sh push qwen:0.5b
```

It will push all the layers, config and the manifest to Codex. A CID for the manifest will be printed at the end.

You can clean the Ollama dir now - e.g. `/usr/share/ollama/modles` and check the model is no longer available

```
sudo rm -rf /usr/share/ollama/.ollama/*
ollama list
```

Try to pull from Codex - the CID of the model will be searched in a [model registry]. If it fails, you can provide the CID in an option `--manifest-cid=${CID}`

```
./collama.sh pull qwen:0.5b
```

You should be able to list the `qwen` model again

```
ollama list
```

## Ollama Codex Model Registry

Since we need the CID to be able to pull the manifest from Codex, but there is no simple and straightforward way to have a decentralized key/value store maintaining the mapping of `model name -> CID`, we go around this by hosting a static "[model registry](https://app.radicle.xyz/nodes/seed.radicle.garden/rad:z4UyZVeb2hy7oKxZ1mUgR1ujAzFsy)" on [Radicle](https://radicle.xyz). 

The registry manifest is automatically updated (downloaded) if the registry file is older than an hour. You can customize the registry URL by various config options (Radicle API URL, repository ID, git ref, filename) and hence easily provide/use your own registry.