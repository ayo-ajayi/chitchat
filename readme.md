## ChitChat  
ChitChat is a command-line conversational AI tool that allows you to use [openai](https://platform.openai.com/) gpt models from the comfort of your terminal. You can work with latest models switching between several of them convieniently. You can also save your conversations and history for later use. Built with [Cobra](https://github.com/spf13/cobra)


### Installation
Ensure you have [Go](https://go.dev/) installed and $GOBIN in your $PATH variable then run the following commands:


```bash 
mv path/to/.chitchat.yaml ~/
go install github.com/ayo-ajayi/chitchat@latest

chitchat -v
```

Run CLI using a .chitchat.yaml file that is not in the home directory:
```bash
chitchat --config path/to/.chitchat.yaml
```
The template for the .chitchat.yaml file can be found [here](./.chitchat.yaml)  
Also ensure that you have [redis](https://redis.io/) running on your machine or on the cloud for storage.

### Getting Started
To use ChitChat, simply run the following command to see the available options:
```bash
chitchat -h
 
chitchat help
```


###  Author
-   [Ayomide Ajayi](https://github.com/ayo-ajayi)
