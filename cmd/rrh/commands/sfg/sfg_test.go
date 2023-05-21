package sfg

import "os"

func Example_Init() {
	cmd := New()
	cmd.SetArgs([]string{"--shell", "zsh"})
	cmd.SetOut(os.Stdout)
	cmd.Execute()
	// Output:
	// cdrrh(){
	//     to_path=$(rrh repository info --entry path $1)
	//     if [[ $? -eq 0 ]]; then
	//         if [[ $(echo $to_path | wc -l) -ne 1 ]]; then
	//             echo "Error: multiple paths are given."
	//             return 1
	//         fi
	//         cd ${to_path#Path: }
	//         pwd
	//     else
	//         return 1
	//     fi
	// }
	//
	// #compdef cdrrh
	// _cdrrh() {
	//     local -a ids
	//     ids=($(rrh repository list --entry id | sort -u))
	//     _values $state $ids
	// }
	// compdef _cdrrh cdrrh
	// rrhpeco(){
	//     csv=$(rrh list --format table | peco)
	//     if [[ $(echo $csv | wc -l) -ne 1 ]]; then
	//         echo "multiple entries are given"
	//         return 1
	//     fi
	//     cd $(echo $csv | awk -F | '{ print $4 }' | tr -d ' ')
	//     pwd
	// }
	// rrhfzf(){
	//     csv=$(rrh list --format table | fzf)
	//     if [[ $(echo $csv | wc -l) -ne 1 ]]; then
	//         echo "multiple entries are given"
	//         return 1
	//     fi
	//     cd $(echo $csv | awk -F | '{ print $4 }' | tr -d ' ')
	//     pwd
	// }
}
