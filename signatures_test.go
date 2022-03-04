package client_test

import client "aurora/client/go"

var (
	addTask0, addTask1, addTask2                      client.Signature
	multiplyTask0, multiplyTask1                      client.Signature
	sumIntsTask, sumFloatsTask, concatTask, splitTask client.Signature
	panicTask                                         client.Signature
	longRunningTask                                   client.Signature
)

var initTasks = func() {
	addTask0 = client.Signature{
		Name: "add",
		Args: []client.Arg{
			{
				Type:  "int64",
				Value: 1,
			},
			{
				Type:  "int64",
				Value: 1,
			},
		},
	}

	addTask1 = client.Signature{
		Name: "add",
		Args: []client.Arg{
			{
				Type:  "int64",
				Value: 2,
			},
			{
				Type:  "int64",
				Value: 2,
			},
		},
	}

	addTask2 = client.Signature{
		Name: "add",
		Args: []client.Arg{
			{
				Type:  "int64",
				Value: 5,
			},
			{
				Type:  "int64",
				Value: 6,
			},
		},
	}

	multiplyTask0 = client.Signature{
		Name: "multiply",
		Args: []client.Arg{
			{
				Type:  "int64",
				Value: 4,
			},
		},
	}

	multiplyTask1 = client.Signature{
		Name: "multiply",
	}

	sumIntsTask = client.Signature{
		Name: "sum_ints",
		Args: []client.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1, 2},
			},
		},
	}

	sumFloatsTask = client.Signature{
		Name: "sum_floats",
		Args: []client.Arg{
			{
				Type:  "[]float64",
				Value: []float64{1.5, 2.7},
			},
		},
	}

	concatTask = client.Signature{
		Name: "concat",
		Args: []client.Arg{
			{
				Type:  "[]string",
				Value: []string{"foo", "bar"},
			},
		},
	}

	splitTask = client.Signature{
		Name: "split",
		Args: []client.Arg{
			{
				Type:  "string",
				Value: "foo",
			},
		},
	}

	panicTask = client.Signature{
		Name: "panic_task",
	}

	longRunningTask = client.Signature{
		Name: "long_running_task",
	}
}
