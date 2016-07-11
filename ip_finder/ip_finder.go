package ip_finder

type IPFinder interface {
	GetIP(cid string) string
}
