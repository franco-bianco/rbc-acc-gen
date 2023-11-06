##### RBC GEN

- Ensure `registered.txt` and `unregistered.txt` files are in the `data` directory
- Add Gmail App Password to Line 29 of Main.go File.

```
	if err := imapClient.Login("testEmail@gmail.com", "asdf asdp sfdg asda"); err != nil {
		log.Fatalf("error logging into IMAP client: %s", err)
	}
```

- To Run:

```
go run ./cmd
```
