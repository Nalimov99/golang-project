package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"garagesale/internal/platform/auth"
	"garagesale/internal/platform/database"
	"garagesale/internal/platform/user"
	"garagesale/internal/schema"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log := log.New(os.Stdout, "ADMIN: ", log.LstdFlags|log.Lshortfile)

	var cfg struct {
		DB database.Config
	}
	if err := envconfig.Process("garagesale", &cfg); err != nil {
		return errors.Wrap(err, "generating config usage")
	}

	var err error

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		err = migrate(cfg.DB)
		if err == nil {
			log.Print("migrations complete")
		}
	case "seed":
		err = seed(cfg.DB)
		if err == nil {
			log.Print("seed complete")
		}
	case "useradd":
		role := flag.Arg(1)
		err = useradd(cfg.DB, role)
		if err == nil {
			log.Print("user added")
		}
	case "keygen":
		err = keygen()
		if err == nil {
			log.Print("keygen success")
		}
	default:
		log.Print("No args passed")
		return nil
	}

	return err
}

func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return errors.Wrap(err, "applying migrations")
	}

	return nil
}

func seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Seed(db); err != nil {
		return errors.Wrap(err, "applying seed")
	}

	return nil
}

func useradd(cfg database.Config, roleFlag string) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	authRoles := []string{}
	switch roleFlag {
	case "user":
		authRoles = append(authRoles, auth.RoleUser)
	default:
		authRoles = append(authRoles, []string{auth.RoleAdmin, auth.RoleUser}...)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSuffix(name, "\n")

	fmt.Print("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	email = strings.TrimSuffix(email, "\n")

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Print("Repeat Password: ")
	repeatBytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	fmt.Println()
	if bytes.Compare(repeatBytePassword, bytePassword) != 0 {
		return errors.New("passwords are not equal")
	}

	nu := user.NewUser{
		Name:            name,
		Email:           email,
		Roles:           authRoles,
		Password:        string(bytePassword),
		PasswordConfirm: string(repeatBytePassword),
	}

	if _, err := user.Create(context.Background(), db, nu, time.Now()); err != nil {
		return err
	}

	return nil
}

func keygen() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	file, err := os.Create("private.pem")
	defer file.Close()
	if err != nil {
		return err
	}

	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err := pem.Encode(file, &block); err != nil {
		return err
	}

	return nil
}
