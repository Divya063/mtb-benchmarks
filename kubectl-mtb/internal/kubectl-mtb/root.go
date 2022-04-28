/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package kubectl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/pkg/benchmark"
)

var rootCmd *cobra.Command
var maxProfileLevel = 3
var benchmarks []*benchmark.Benchmark

// singular-plural
var supportedResourceNames = sets.NewString("benchmarks", "benchmark")

func init() {

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "kubectl-mtb",
		Short: "Multi-Tenancy Benchmarks",
	}

	rootCmd.PersistentFlags().IntP("profile-level", "p", maxProfileLevel, "ProfileLevel of the benchmarks.")

	// Commands
	rootCmd.AddCommand(newGetCmd())
	rootCmd.AddCommand(newRunCmd())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
