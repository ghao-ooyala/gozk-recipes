package recipes

import (
	. "launchpad.net/gocheck"
	"sort"
)

func (s *RecipesSuite) TestFunctionsOnInvalidNodes(c *C) {
	zk.RecursiveDelete("/test")
	_, _, err := zk.Get("/test/hi")
	c.Check(err, Not(IsNil))
	stat, _ := zk.Exists("/test/hi")
	c.Check(stat, IsNil)
	err = zk.Delete("/test/hi", -1)
	c.Check(err, Not(IsNil))
	err = zk.RecursiveDelete("/test/hi")
	c.Check(err, Not(IsNil))
	nodes, _, err := zk.Children("/test/hi")
	c.Check(nodes, DeepEquals, []string(nil))
	c.Check(err, Not(IsNil))
	nodes, _, err = zk.VisibleChildren("/test/hi")
	c.Check(nodes, DeepEquals, []string(nil))
	c.Check(err, Not(IsNil))
}

func (s *RecipesSuite) TestFunctionsOnValidNodes(c *C) {
	zk.RecursiveDelete("/test")
	// test TouchAndSet, Touch
	_, err := zk.TouchAndSet("/test/test", "omgwtfbbq")
	c.Check(err, IsNil)
	path, err := zk.Touch("/test/test|dir") // test pipe in name
	c.Check(err, IsNil)
	c.Check(path, Equals, "/test/test|dir")
	_, err = zk.TouchAndSet("/test/test|dir/subdir", "dirwithdata")
	c.Check(err, IsNil)
	_, err = zk.TouchAndSet("/test/test|dir/subdir/ooga", "booga")
	c.Check(err, IsNil)
	_, err = zk.TouchAndSet("/test/test|dir/subdir/booga", "ooga")
	c.Check(err, IsNil)
	_, err = zk.TouchAndSet("/test/test|dir/subdir/omgwtf", "bbq")
	c.Check(err, IsNil)
	path, err = zk.Touch("/test/test:hidden") // test colon in name
	c.Check(err, IsNil)
	c.Check(path, Equals, "/test/test:hidden")
	_, err = zk.TouchAndSet("/test/test:hidden/.hidden", "this isn't the hidden node you're looking for")
	c.Check(err, IsNil)
	_, err = zk.TouchAndSet("/test/test:hidden/.alsohidden", "move along")
	c.Check(err, IsNil)
	_, err = zk.TouchAndSet("/test/test:hidden/nothidden", "damnit you can see me")
	c.Check(err, IsNil)
	_, err = zk.TouchAndSet("/test/test:hidden/alsonothidden", "and me too")
	c.Check(err, IsNil)

	// test Exists
	stat, _ := zk.Exists("/test/test")
	c.Check(stat, Not(IsNil))
	stat, _ = zk.Exists("/test/nothere")
	c.Check(stat, IsNil)
	stat, _ = zk.Exists("/test/test|dir")
	c.Check(stat, Not(IsNil))
	stat, _ = zk.Exists("/test/test|dir/subdir")
	c.Check(stat, Not(IsNil))
	stat, _ = zk.Exists("/test/test|dir/subdir/ooga")
	c.Check(stat, Not(IsNil))
	stat, _ = zk.Exists("/test/test:hidden")
	c.Check(stat, Not(IsNil))
	stat, _ = zk.Exists("/test/test:hidden/.hidden")
	c.Check(stat, Not(IsNil))
	stat, _ = zk.Exists("/test/test:hidden/nothidden")
	c.Check(stat, Not(IsNil))

	// test Get
	data, _, err := zk.Get("/test/test")
	c.Check(data, Equals, "omgwtfbbq")
	c.Check(err, IsNil)
	data, _, err = zk.Get("/test/test|dir")
	c.Check(data, Equals, "")
	c.Check(err, IsNil)
	data, _, err = zk.Get("/test/test|dir/subdir")
	c.Check(data, Equals, "dirwithdata")
	c.Check(err, IsNil)
	data, _, err = zk.Get("/test/test|dir/subdir/booga")
	c.Check(data, Equals, "ooga")
	c.Check(err, IsNil)
	data, _, err = zk.Get("/test/test:hidden")
	c.Check(data, Equals, "")
	c.Check(err, IsNil)
	data, _, err = zk.Get("/test/test:hidden/.alsohidden")
	c.Check(data, Equals, "move along")
	c.Check(err, IsNil)
	data, _, err = zk.Get("/test/test:hidden/alsonothidden")
	c.Check(data, Equals, "and me too")
	c.Check(err, IsNil)

	// test Children
	nodes, _, err := zk.Children("/test/test")
	c.Check(err, IsNil)
	c.Check(nodes, DeepEquals, []string(nil))
	nodes, _, err = zk.Children("/test/test|dir")
	c.Check(err, IsNil)
	c.Check(nodes, DeepEquals, []string{"subdir"})
	nodes, _, err = zk.Children("/test/test|dir/subdir")
	c.Check(err, IsNil)
	sort.Strings(nodes)
	c.Check(nodes, DeepEquals, []string{"booga", "omgwtf", "ooga"})
	nodes, _, err = zk.Children("/test/test:hidden")
	c.Check(err, IsNil)
	sort.Strings(nodes)
	c.Check(nodes, DeepEquals, []string{".alsohidden", ".hidden", "alsonothidden", "nothidden"})

	// test FilterHidden, VisibleChildren
	nodes, _, err = zk.VisibleChildren("/test/test")
	c.Check(err, IsNil)
	c.Check(nodes, DeepEquals, []string(nil))
	nodes, _, err = zk.VisibleChildren("/test/test|dir")
	c.Check(err, IsNil)
	c.Check(nodes, DeepEquals, []string{"subdir"})
	nodes, _, err = zk.VisibleChildren("/test/test|dir/subdir")
	c.Check(err, IsNil)
	sort.Strings(nodes)
	c.Check(nodes, DeepEquals, []string{"booga", "omgwtf", "ooga"})
	nodes, _, err = zk.VisibleChildren("/test/test:hidden")
	c.Check(err, IsNil)
	sort.Strings(nodes)
	c.Check(nodes, DeepEquals, []string{"alsonothidden", "nothidden"})

	// test Delete
	err = zk.Delete("/test/test", -1)
	c.Check(err, IsNil)
	stat, _ = zk.Exists("/test/test")
	c.Check(stat, IsNil)
	err = zk.Delete("/test/test|dir", -1)
	c.Check(err, Not(IsNil))
	stat, _ = zk.Exists("/test/test|dir")
	c.Check(stat, Not(IsNil))
	err = zk.Delete("/test/test:hidden", -1)
	c.Check(err, Not(IsNil))
	stat, _ = zk.Exists("/test/test:hidden")
	c.Check(stat, Not(IsNil))

	// test RecursiveDelete
	err = zk.RecursiveDelete("/test/test|dir")
	c.Check(err, IsNil)
	stat, _ = zk.Exists("/test/test|dir")
	c.Check(stat, IsNil)
	stat, _ = zk.Exists("/test/test|dir/subdir")
	c.Check(stat, IsNil)
	stat, _ = zk.Exists("/test/test|dir/subdir/ooga")
	c.Check(stat, IsNil)
	err = zk.RecursiveDelete("/test/test:hidden")
	c.Check(err, IsNil)
	stat, _ = zk.Exists("/test/test:hidden")
	c.Check(stat, IsNil)
	stat, _ = zk.Exists("/test/test:hidden/.hidden")
	c.Check(stat, IsNil)
}