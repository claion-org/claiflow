# MultiStep Features

This example demonstrates how user-registered commands can be executed in multiple ways using the basic features of Multistep.

## Features
- Simple: execute a single step with a command.
- Iteration: execute the iteration step.
- Passing a value: passing values between multiple commands.

## How to run
1. Register your own functions as commands. In this example, we'll register three commands.
    ```go
    // return "hello <name>"
	if err := c.RegisterCommand("helloworld", HelloWorld); err != nil {
		log.Fatal(err)
	}

    // return x**y
	if err := c.RegisterCommand("math_pow", MathPow); err != nil {
		log.Fatal(err)
	}

    // return value1, value2 -> value2, value1
	if err := c.RegisterCommand("swap_command", SwapCommand); err != nil {
		log.Fatal(err)
	}
    ```

2. Registered commands must exist on the server as a Template in order to be used.

    The Template has a flow field that specifies how the command should be executed.

    This example shows three different ways: simple execution, iteration execution, and passing a value.

    Insert the following examples into your server database.

    * Simple
        ```sql
        REPLACE INTO `template` (`uuid`, `name`, `summary`, `flow`, `inputs`, `origin`, `created`) VALUES ('example_simple', 'example_simple', '', '[{"$id":"step1","$command":"helloworld","inputs":"$inputs"}]', '{}', 'userdefined', NOW());
        ```

    * Iteration
        ```sql
        REPLACE INTO `template` (`uuid`, `name`, `summary`, `flow`, `inputs`, `origin`, `created`) VALUES ('example_iter', 'example_iter', '', '[{"$id":"step1","$range":"$inputs.x_list","$steps":[{"$id":"step2","$command":"math_pow","inputs":{"x":"$step1.val","y":2}}]}]', '{}', 'userdefined', NOW());
        ```

    * Passing a value
        ```sql
        REPLACE INTO `template` (`uuid`, `name`, `summary`, `flow`, `inputs`, `origin`, `created`) VALUES ('example_pass_val', 'example_pass_val', '', '[{"$id":"step1","$command":"swap_command","inputs":{"param1":"$inputs.input1","param2":"$inputs.input2"}},{"$id":"step2","$command":"swap_command","inputs":{"param1":"$step1.outputs.value1","param2":"$step1.outputs.value2"}}]', '{}', 'userdefined', NOW());
        ```

3. Run client.
    ```bash
    go run main.go \
        --server=localhost:18099 \
        --clusteruuid=<CLUSTER_UUID> \
        --token=<CLUSTER_TOKEN>
    ```

4. Request the server to create a Service with the Template you created above.
    * Simple
        ```bash
        curl -X POST http://localhost:8099/api/v1/service \
            -H 'Content-Type: application/json' \
            -d '{"cluster_uuid":"<CLUSTER_UUID>","name":"example_simple","template_uuid":"example_simple","inputs":{"name":"world"}}'
        ```
        
    * Iteration
        ```bash
        curl -X POST http://localhost:8099/api/v1/service \
            -H 'Content-Type: application/json' \
            -d '{"cluster_uuid":"<CLUSTER_UUID>","name":"example_iter","template_uuid":"example_iter","inputs":{"x_list":[1,2,3,4,5]}}'
        ```
    
    * Passing a value
        ```bash
        curl -X POST http://localhost:8099/api/v1/service \
            -H 'Content-Type: application/json' \
            -d '{"cluster_uuid":"<CLUSTER_UUID>","name":"example_pass_val","template_uuid":"example_pass_val","inputs":{"input1":"foo","input2":"bar"}}'
        ```