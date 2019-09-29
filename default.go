package log

var DefaultProvider = New()

var Fatal = DefaultProvider.V(FATAL).Str
var Fatalf = DefaultProvider.V(FATAL).Strf
var Error = DefaultProvider.V(ERROR).Str
var Errorf = DefaultProvider.V(ERROR).Strf
var Warning = DefaultProvider.V(WARNING).Str
var Warningf = DefaultProvider.V(WARNING).Strf
var Info = DefaultProvider.V(INFO).Str
var Infof = DefaultProvider.V(INFO).Strf
