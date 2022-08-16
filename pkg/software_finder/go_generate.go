//go:build ignore

package software_finder

//go:generate mockgen -source=software_finder.go -destination=software_finder_mock_test.go -package=$GOPACKAGE -write_package_comment=false
