LOCKS_OFF_TESTS_DIV_ID = 'locks_off_tests';


/**
 @param {list} data_list --- Each element is a list containing
 integers and floats.  Indices are formatted according to 
 */
function draw_locks_off_graph(data_list)
{
    // keys are number of operations performed, values are lists of
    // floats.  Each float is measured throughput of that condition.
    var locks_off_graph_data_dict = {};
    // each element is the number of reads performed on a particular
    // condition.
    var num_threads_list = [];
    for (var index in data_list)
    {
        var condition_data = new RunStats(data_list[index]);
        if (! (condition_data.num_threads in locks_off_graph_data_dict))
        {
            locks_off_graph_data_dict[condition_data.num_threads] = [];
            num_threads_list.push(condition_data.num_threads);
        }
        locks_off_graph_data_dict[condition_data.num_threads].push(condition_data);
    }

    // sort num_threads_list so that graph conditions will be displayed
    // in order
    num_threads_list.sort();

    // flatten data into a single average.  Each element in this list
    // is a RunStats object
    var average_data_list = [];
    
    for (index in locks_off_graph_data_dict)
    {
        var avg_stats = stats_avg_throughput(locks_off_graph_data_dict[index]);
        average_data_list.push(avg_stats);
    }

    // sort warm_graph_average_data_list by num_reads so that plot
    // data in order.
    average_data_list.sort(
        function (a,b)
        {
            if (a.num_threads < b.num_threads)
                return -1;
            if (a.num_threads > b.num_threads)
                return 1;
            return 0;
        });
    
    // Handle plotting averaged data
    plot(LOCKS_OFF_TESTS_DIV_ID,average_data_list);
}

