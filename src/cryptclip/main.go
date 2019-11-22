package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/atotto/clipboard"
)

func main() {
	timeoutPtr := flag.Duration("t", 0, "Erase clipboard after timeout.  Durations are specified like \"20s\" or \"2h45m\".  0 (default) means never erase.")
	encryptPtr := flag.Bool("enc", false, "encrypt the provided string")
	decryptPtr := flag.Bool("dec", false, "decrypt the provided string")
	keyPtr := flag.String("key", "", "The key to use for encryption")
	msgPtr := flag.String("msk", "", "The message to encrypt/decrypt")

	flag.Parse()

	if *encryptPtr && *decryptPtr {
		os.Exit(1)
	}

	if *encryptPtr {
		out, err := encrypt([]byte(*keyPtr), *msgPtr)
		if err != nil {
			log.Fatal(err)
		}

		err = clipboard.WriteAll(string(out))
		if err != nil {
			log.Fatal(err)
		}

		if timeoutPtr != nil && *timeoutPtr > 0 {
			<-time.After(*timeoutPtr)
			text, err := clipboard.ReadAll()
			if err != nil {
				os.Exit(0)
			}

			if text == string(out) {
				err = clipboard.WriteAll("")
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if err != nil {
			os.Exit(1)
		}
	} else if *decryptPtr {
		out, err := decrypt([]byte(*keyPtr), *msgPtr)
		if err != nil {
			log.Fatal(err)
		}

		err = clipboard.WriteAll(string(out))
		if err != nil {
			log.Fatal(err)
		}

		if timeoutPtr != nil && *timeoutPtr > 0 {
			<-time.After(*timeoutPtr)
			text, err := clipboard.ReadAll()
			if err != nil {
				os.Exit(0)
			}

			if text == string(out) {
				err = clipboard.WriteAll("")
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if err != nil {
			os.Exit(1)
		}
	}
}
