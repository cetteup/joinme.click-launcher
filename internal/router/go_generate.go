//go:build ignore

package router

//go:generate mockgen -source=router.go -destination=router_mock_test.go -package=$GOPACKAGE -write_package_comment=false
