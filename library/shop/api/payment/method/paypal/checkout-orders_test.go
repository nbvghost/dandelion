package paypal

import "testing"

//var matchRegexp = regexp.MustCompile("^(https:)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png|jpeg|JPG|GIF|PNG|JPEG)")
func Test_matchImageUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "match image url https", args: args{url: "https://www.sdfds.ssd/sdfds/sdf/sdf/dsf/sd/fds.png"}, want: true},
		{name: "match image url http", args: args{url: "http://www.sdfds.ssd/sdfds/sdf/sdf/dsf/sd/fds.png"}, want: false},
		{name: "match image url webp", args: args{url: "http://www.sdfds.ssd/sdfds/sdf/sdf/dsf/sd/fds.webp"}, want: false},
		{name: "match image url jpeg", args: args{url: "https://www.sdfds.ssd/sdfds/sdf/sdf/dsf/sd/fds.jpeg"}, want: true},
		{name: "match image url JPG", args: args{url: "https://www.sdfds.ssd/sdfds/sdf/sdf/dsf/sd/fds.JPG"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchImageUrl(tt.args.url); got != tt.want {
				t.Errorf("matchImageUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
