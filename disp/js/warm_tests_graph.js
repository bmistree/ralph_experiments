WARM_TESTS_DIV_ID = 'warm_tests';


/**
 @param {list} data_list --- Each element is a list containing
 integers and floats.  Indices are formatted according to 
 */
function draw_read_warm_graph(data_list)
{
    // keys are number of operations performed, values are lists of
    // floats.  Each float is measured throughput of that condition.
    var warm_graph_data_dict = {};
    // each element is the number of reads performed on a particular
    // condition.
    var num_reads_list = [];
    for (var index in data_list)
    {
        var condition_data = new RunStats(data_list[index]);
        if (! (condition_data.num_reads in warm_graph_data_dict))
        {
            warm_graph_data_dict[condition_data.num_reads] = [];
            num_reads_list.push(condition_data.num_reads);
        }
        warm_graph_data_dict[condition_data.num_reads].push(condition_data);
    }

    // sort num_reads_list so that graph conditions will be displayed
    // in order
    num_reads_list.sort();
    
    // flatten data into a single average.  Each element in this list
    // is a RunStats object
    var warm_graph_average_data_list = [];
    
    for (index in warm_graph_data_dict)
    {
        var avg_stats = stats_avg_throughput(warm_graph_data_dict[index]);
        warm_graph_average_data_list.push(avg_stats);
    }

    // sort warm_graph_average_data_list by num_reads so that plot
    // data in order.
    warm_graph_average_data_list.sort(
        function (a,b)
        {
            if (a.num_reads < b.num_reads)
                return -1;
            if (a.num_reads > b.num_reads)
                return 1;
            return 0;
        });
    

    // Handle plotting averaged data
    plot(WARM_TESTS_DIV_ID,warm_graph_average_data_list);
}
