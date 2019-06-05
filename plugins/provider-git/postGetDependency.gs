// Provider : name of provider (npm)
// Name : name of downloaded package
// Version : version of downloaded package
// Url : URL of downloaded package
// Path  : Future path of downloaded package
// Result : binary downloaded

var folder = unzip(Result);

var firstChildrenName = Object.keys(folder.Children)[0];
var firstChildren = folder.Children[firstChildrenName];

saveFileAt(firstChildren, Path);
saveLockDep(Path);

Path;
