package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/atotto/clipboard"
)

var (
	timeoutPtr = flag.Duration("t", 0, "Erase clipboard after timeout.  Durations are specified like \"20s\" or \"2h45m\".  0 (default) means never erase.")
	encryptPtr = flag.Bool("e", false, "encrypt the provided string")
	decryptPtr = flag.Bool("d", false, "decrypt the provided string")
	keyPtr     = flag.String("key", "", "The key to use for encryption")
	msgPtr     = flag.String("msg", "", "The message to encrypt/decrypt")
)

func main() {
	flag.Parse()

	if *encryptPtr && *decryptPtr {
		os.Exit(1)
	}

	var key string
	var msg string

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if keyPtr == nil || len(*keyPtr) <= 0 {
		tmpKey, err := ioutil.ReadFile(path.Join(usr.HomeDir, ".keys/clipcrypt.key"))
		if err != nil {
			log.Fatal(err)
		}

		key = string(tmpKey)
	} else {
		key = *keyPtr
	}

	if msgPtr == nil || len(*msgPtr) <= 0 {
		tmpMsg, err := clipboard.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		msg = tmpMsg
	} else {
		msg = *msgPtr
	}

	if *encryptPtr {
		out, err := encrypt([]byte(key), msg)
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
		out, err := decrypt([]byte(key), msg)
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
