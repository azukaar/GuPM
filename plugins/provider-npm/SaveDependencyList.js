var package = readJsonFile('./package.json')
package.dependencies = {};

for(d in Dependencies) {
    var dep = Dependencies[d];
    package.dependencies[dep.name] = dep.version;
}

saveJsonFile('./package.json', package)