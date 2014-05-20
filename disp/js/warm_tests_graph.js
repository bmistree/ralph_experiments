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
        if (! (condition_data.reads_per_second in warm_graph_data_dict))
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
    var num_conditions = warm_graph_average_data_list.length;
    var width =
        (BAR_WIDTH + BAR_CATEGORY_SPACING_WIDTH)*num_conditions;

    
    var bar_chart = d3.select('#' + WARM_TESTS_DIV_ID).
        append('svg:svg').
        attr('width', width).
        attr('height', BAR_HEIGHT);

    var x_rect_positions =
        d3.scale.linear().domain([0, num_conditions]).range([0, width]);
    var y_heights =
        d3.scale.linear().domain(
            [0, d3.max(warm_graph_average_data_list,
                       function(datum) { return datum.reads_per_second; })]).
        rangeRound([0, BAR_HEIGHT]);
    
    // FIXME: don't understand semantics of selectAll
    bar_chart.selectAll('rect').
        data(warm_graph_average_data_list).
        enter().
        append("svg:rect").
        attr("x",
             function(datum, index)
             {
                 return x_rect_positions(index);
             }).
        attr("y",
             function(datum)
             {
                 return BAR_HEIGHT - y_heights(datum.reads_per_second);
             }).
        attr("height",
             function(datum)
             {
                 return y_heights(datum.reads_per_second);
             }).
        attr("width", BAR_WIDTH).
        attr("fill", "steelblue");
}
