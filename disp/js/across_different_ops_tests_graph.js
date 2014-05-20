
ACROSS_DIFFERENT_OPS_TESTS_DIV_ID = 'across_different_ops';

/**
 @param {list} data_list --- Each element is a list containing
 integers and floats.  Indices are formatted according to 
 */
function across_different_ops_graph(data_list)
{
    // keys are number of operations performed, values are lists of
    // floats.  Each float is measured throughput of that condition.
    var graph_data_dict = {};

    for (var index in data_list)
    {
        var condition_data = new RunStats(data_list[index]);
        var new_index = produce_index_across_ops(
            condition_data.num_threads,condition_data.read_type);

        if (! (new_index in graph_data_dict))
            graph_data_dict[new_index] = [];

        graph_data_dict[new_index].push(condition_data);
    }

    // flatten data into a single average.  Each element in this list
    // is a RunStats object
    var average_data_list = [];
    
    for (index in graph_data_dict)
    {
        var avg_stats = stats_avg_throughput(graph_data_dict[index]);
        average_data_list.push(avg_stats);
    }

    // sort warm_graph_average_data_list by num_reads so that plot
    // data in order.
    average_data_list.sort(
        function (a,b)
        {
            if (a.read_type < b.read_type)
                return -1;
            if (a.read_type > b.read_type)
                return 1;

            if (a.num_threads < b.num_threads)
                return -1;
            if (a.num_threads > b.num_threads)
                return 1;
            return 0;
        });
    
    // Handle plotting averaged data
    plot(ACROSS_DIFFERENT_OPS_TESTS_DIV_ID,average_data_list);
}

function produce_index_across_ops(num_threads,read_type)
{
    return read_type + "|" + num_threads;
}