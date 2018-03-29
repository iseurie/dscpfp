package main

import (
	"bufio"
	"flag"
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"golang.org/x/crypto/ssh/terminal"
	"hash/crc32"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func errcode(what string, cause error) int {
	m := what + fmt.Sprintf("%v", cause)
	return (int)(crc32.ChecksumIEEE([]byte(m)))
}

func errck(what string, cause error) {
	if cause != nil {
		if what != "" {
			fmt.Print(what + ": ")
		}
		fmt.Println(cause)
		os.Exit(errcode(what, cause))
	}
}

func uid_search(s *dg.Session, tag string) (uid string, err error) {
	spl := strings.SplitN(tag, "#", 2)
	uname := spl[0]
	discriminator := ""
	if len(spl) > 1 {
		discriminator = spl[1]
	}
	relationships, err := s.RelationshipsGet()
	if err != nil {
		return "", fmt.Errorf("Retrieve contacts: %s", err)
	}

	for _, r := range relationships {
		d := discriminator == "" || discriminator == r.User.Discriminator
		if d && strings.Contains(r.User.Username, uname) {
			return r.User.ID, nil
		}
	}
	return "", fmt.Errorf("No matching tag found")
}

func credentials_or_die() (uname string, pass string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("email: ")
	uname, err := reader.ReadString('\n')
	errck("", err)
	fmt.Print("pass: ")
	rawpass, err := terminal.ReadPassword(int(syscall.Stdin))
	errck("", err)
	pass = strings.TrimSpace(string(rawpass))
	uname = strings.TrimSpace(uname)
	fmt.Println()
	return
}

var szp uint
var uid string
var opath string
var token string
var uname string

func main() {
	flag.UintVar(&szp, "dimen", 10, "Avatar asset will be retrieved for dimension 1 << this.")
	flag.StringVar(&opath, "opath", "", "Output path. If set, profile image will be retrieved and output to this file. Extensions will be appended for the appropriate MIME type.")
	flag.StringVar(&uid, "uid", "", "User ID for which to fetch avatar.")
	flag.StringVar(&uname, "uname", "", "Username and optional tag for which to fetch avatar.")
	flag.StringVar(&token, "token", "", "Bot token for Discord API.")
	flag.Parse()
	if uid == "" && uname == "" {
		fmt.Fprintf(os.Stderr, "Please provide -uid or -uname.\n")
		os.Exit(1)
	}
	if token == "" {
		token = os.Getenv("DSCPFP_TOKEN")
	}
	if token == "" {
		fmt.Fprintf(os.Stderr, "Please provide a token either via the command-line or `DSCPFP_TOKEN`.\n")
		os.Exit(1)
	}
	session, err := dg.New("Bot " + token)
	errck("Session auth failure", err)
	if uid == "" {
		// uname must be set
		u, p := credentials_or_die()
		fmt.Println("Logging...")
		searchSession, err := dg.New(u, p)
		errck("Authenticate session for `uname` search", err)
		uid, err = uid_search(searchSession, uname)
		errck("", err)
	}
	u, err := session.User(uid)
	errck("Retrieve target user", err)
	aUri := u.AvatarURL(strconv.FormatUint((uint64)(0x01<<szp), 10))
	if opath == "" {
		fmt.Println(aUri)
		os.Exit(0)
	}
	resp, err := http.Get(aUri)
	errck(fmt.Sprintf("HTTP/GET `%s`", aUri), err)
	defer resp.Body.Close()
	exts, _ := mime.ExtensionsByType(resp.Header.Get("Content-Type"))
	ext := filepath.Ext(opath)
	if len(exts) > 0 {
		ext = exts[0]
	} else {
		fmt.Fprintf(os.Stderr, "Warning: No extension found for provided Content-Type `%s`; not appending\n", exts)
	}
	if filepath.Ext(opath) != ext {
		opath += ext
	}
	ostrm, err := os.OpenFile(opath, os.O_WRONLY|os.O_CREATE, 0555)
	errck("output file", err)
	_, err = io.Copy(ostrm, resp.Body)
	errck("output file", err)
}
