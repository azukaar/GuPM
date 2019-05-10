var goMod = readFile(Dependency.path + '/go.mod');

if(goMod) {
    var requires = goMod.match(/require\s\(\n?([\s\w\_\-\.\/\n]+)\n?\)/)[1].split(/\n/);
    Dependency.dependencies = [];
    
    for(r in requires) {
        var require = requires[r];
        if(require != "") {
            if(require.split(" ").length>1) {
                Dependency.dependencies.push({
                    provider: 'go',
                    name: require.split(" ")[0].trim(),
                    version: require.split(" ")[1].trim()
                })
            }
        }
    }
}

Dependency;
