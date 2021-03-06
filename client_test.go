package nipap
import (
	"testing"
	"net"
	sq "github.com/mrevilme/go-nipap/search_query"
)

func Test_CreateClient(t *testing.T) {
	NewTestClient(t)
}

func Test_Call(t *testing.T) {
	_,c := NewClient("http://localhost:1337/RPC2","svc_oms@local","svc_oms")
	x := c.Run("echo", map[string]interface{}{"message":"foobar"}, nil)
	t.Fatal(x)
}

func Test_List_All_Prefixes(t *testing.T) {
	c := NewTestClient(t)
	err, prefixes := c.ListPrefix(nil)
	if err != nil {
		t.Fatal(err)
	} else {
		if len(prefixes) <= 1 {
			t.Fatalf("Should get more then 1 prefix, got %d prefixes", len(prefixes))
		}
	}
}

func Test_List_Specific_Prefix(t *testing.T) {
	c := NewTestClient(t)
	err, prefixes := c.ListPrefix(map[string]string{"prefix":"172.16.5.0/24"})
	if err != nil {
		t.Fatal(err)
	} else {
		if len(prefixes) == 1 {
			p := prefixes[0]
			if p.Prefix != "172.16.5.0/24" {
				t.Fatalf("Should have recived a prefix 172.16.5.0/24 but recived %s", p.Prefix)
			}
		} else {
			t.Fatalf("Should get more then 1 prefix, got %d prefixes", len(prefixes))
		}
	}

}

func Test_Add_Prefix(t *testing.T) {
	c := NewTestClient(t)
	p := Prefix{}
	p.Prefix = "172.16.4.0/29"
	p.Description = "Hello world"
	p.Type = PrefixTypeReservation
	err, p := c.AddPrefix(p)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Add_Prefix_From_Prefix(t *testing.T) {
	c := NewTestClient(t)
	newP := Prefix{}
	newP.Description = "Prefix from Prefix"
	newP.Type = PrefixTypeReservation

	oldP := Prefix{}
	oldP.Prefix = "172.16.5.0/24"
	err, p := c.AddPrefixFromPrefix(newP, oldP, 29)
	if err != nil {
		t.Fatal(err)
	}
	_,cNet,err := net.ParseCIDR(p.Prefix)
	_,pNet,err := net.ParseCIDR(oldP.Prefix)
	if !pNet.Contains(cNet.IP) {
		t.Fatalf("%s is not in the parent prefix %s", cNet.String(),pNet.String())
	}

}

func Test_Smart_Search_prefix(t *testing.T) {
	c := NewTestClient(t)
	err, prefixes := c.PrefixSmartSearch("172.16.5.1",nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatalf("%+v\n", prefixes)
}

func Test_Smart_Search_Prefix_With_Options(t *testing.T) {
	c := NewTestClient(t)
	o := SearchOptions{}
	o.ParentsDepth = 1
	err, prefixes := c.PrefixSmartSearch("172.16.5.1/32 #foo #bar",&o)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatalf("%+v\n", prefixes)
}

func Test_Search_Prefix(t *testing.T) {
	c := NewTestClient(t)
	o := SearchOptions{}
	o.ParentsDepth = 0

	// Query build
	prefix := sq.Contains("prefix","172.16.5.1/32")
	typeSearch := sq.Equals("type","assignment")
	tagCpeSearch := sq.Or(sq.EqualsAny("tags","foobar"),sq.EqualsAny("inherited_tags","foobar"))
	tagOmsSearch := sq.Or(sq.EqualsAny("tags","bar"),sq.EqualsAny("inherited_tags","bar"))
	tagsSearch := sq.And(tagCpeSearch,tagOmsSearch)

	query := sq.And(tagsSearch,sq.And(prefix,typeSearch))
	query = tagCpeSearch

	err, prefixes := c.SearchPrefix(query,&o)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatalf("%+v\n", prefixes)
}

func Test_DeletePrefix(t *testing.T) {
	c := NewTestClient(t)
	prefix := Prefix{}
	prefix.Id = 95
	err := c.DeletePrefix(prefix,false)
	if err != nil {
		t.Fatal(err)
	}
}


func NewTestClient(t *testing.T) (*Client) {
	err,c := NewClient("http://localhost:1337/RPC2","admin@local","foobar")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatalf("Returned client is nil")
	}

	return c
}
