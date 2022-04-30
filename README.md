# Usage

Example:<br>

Constraints<br>
- 1 &le; M &le; N &le; 1000<br>
- 1 &le; i &le; N<br>
- 1 &le; j &le; M<br>
- 1 &le; A<sub>i</sub> &le; 10<sup>9</sup><br>
- 1 &le; B<sub>j</sub> &le; 10<sup>9</sup>
- 1 &le; C<sub>i,j</sub> &le; 10<sup>9</sup>
- All values in input are integers.

Input<br>
<pre>
N&#009;M
A<sub>1</sub>&#009;A<sub>2</sub>&#009;...&#009;A<sub>N</sub>
B<sub>1</sub>&#009;B<sub>2</sub>&#009;...&#009;B<sub>M</sub>
C<sub>1,1</sub>&#009;C<sub>1,2</sub>&#009;...&#009;C<sub>1,M</sub>
C<sub>2,1</sub>&#009;C<sub>2,2</sub>&#009;...&#009;C<sub>2,M</sub>
&#8286;
C<sub>n,1</sub>&#009;C<sub>n,2</sub>&#009;...&#009;C<sub>N,M</sub>
</pre>

Output<br>
<pre>
T<sub>1,1</sub>&#009;T<sub>1,2</sub>&#009;T<sub>1,3</sub>
T<sub>2,2</sub>&#009;T<sub>2,3</sub>&#009;T<sub>2,3</sub>
S
X&#009;Y&#009;Z
</pre>
The values of T<sub>i,j</sub>, X, Y and Z are integers, and the value of S is a string.


1. Set the input and output structure of template.go
    ~~~go
    type (
        input struct {
            N int
            M int
            A []int   `size:"N"`
            B []int   `size:"M"`
            C [][]int `size:"N,M"`
        }
    
        output struct {
            T    [][]int
            S    string
            X, Y int `EOL:"false"`
            Z    int
        }
    )
    ~~~
    also
    ~~~go
    type (
        input struct {
            N, M int
            A, B []int
            C    [][]int
        }
    
        output struct {
            T   [][]int
            S   string
            XYZ []int
        }
    )
    ~~~
    - Supported types are int, uint, float32, float64, string and slices up to 2 dimensions.
    - When inputting and outputting, the fields of each structure are read and written in the order in which they appear.
    - Each field must be exposed by starting with a capital letter.
    - Non-slice fields read the next element separated by a space or newline.
    - Slice fields are read to the next newline by default. You can also specify the size explicitly with the `size` tag. If the slices are nested, specify the value of the `size` tag as a comma-separated list.
    - By default, each field in the output structure is separated by a newline. This can be changed to be separated by a space from the field immediately following by specifying `false` in the `EOL` tag.

1. Coding the solution to the problem in the solve function of template.go.
    ~~~go
    func solve(in input) output {
        Your Solution
    }
    ~~~