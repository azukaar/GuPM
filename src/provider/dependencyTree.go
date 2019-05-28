package provider

func eliminateRedundancy(tree []map[string]interface {}, path map[string]bool) []map[string]interface {} {
	var cleanTree = make([]map[string]interface {}, 0)
	for index, dep := range tree {
		if(dep["name"] != nil) {
			_ = index
			depKey := dep["name"].(string) + "@" + dep["version"].(string)
			if(path[depKey] != true) {
				cleanTree = append(cleanTree, dep)
			}
		}
	}
	
	for index, dep := range cleanTree {
		if(dep["name"] != nil) {
			nextDepList, ok := dep["dependencies"].([]map[string]interface {})

			if(ok) {
				depKey := dep["name"].(string) + "@" + dep["version"].(string)
				newPath := make(map[string]bool)
				for key, value := range path {
					newPath[key] = value
				}
				newPath[depKey] = true
				newSubTree := eliminateRedundancy(nextDepList, newPath)
				cleanTree[index]["dependencies"] = newSubTree
			}
		}
	}
	return cleanTree
}

func flattenDependencyTree(tree []map[string]interface {}, subTree []map[string]interface {}) ([]map[string]interface {}, []map[string]interface {}) {
	var cleanTree = make([]map[string]interface {}, 0)

	for index, dep := range subTree {
		var rootDeps = make(map[string]string)

		for _, dep := range tree {
			rootDeps[dep["name"].(string)] = dep["version"].(string)
		}

		if(rootDeps[dep["name"].(string)] == "") {
			tree = append(tree, dep)

			nextDepList, ok := dep["dependencies"].([]map[string]interface {})
	
			if(ok) {
				newTree, newSubTree := flattenDependencyTree(tree, nextDepList)
				tree = newTree
				subTree[index]["dependencies"] = newSubTree
			}
		} else if(rootDeps[dep["name"].(string)] != dep["version"].(string)) {
			nextDepList, ok := dep["dependencies"].([]map[string]interface {})
	
			if(ok) {
				newTree, newSubTree := flattenDependencyTree(tree, nextDepList)
				tree = newTree
				subTree[index]["dependencies"] = newSubTree
			}

			cleanTree = append(cleanTree, subTree[index])
		}
	}

	return tree, cleanTree
}

func BuildDependencyTree(tree []map[string]interface {}) []map[string]interface {} {
	cleanTree := eliminateRedundancy(tree, make(map[string]bool))

	for index, dep := range cleanTree {
		nextDepList, ok := dep["dependencies"].([]map[string]interface {})

		if(ok) {
			newCleanTree, newDepList := flattenDependencyTree(cleanTree, nextDepList)
			cleanTree = newCleanTree
			cleanTree[index]["dependencies"] = newDepList
		}
	}
	return cleanTree
}
