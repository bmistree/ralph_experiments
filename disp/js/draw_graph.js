NUM_READS_INDEX = 0;
NUM_THREADS_INDEX = 1;
READ_TYPE_INDEX = 2;
READS_PER_SECOND_INDEX = 3;

WARM_TESTS_DIV_ID = 'warm_tests';
WARM_TESTS_NUM_CONDITIONS = 6;

BAR_WIDTH = 40;
BAR_CATEGORY_SPACING_WIDTH = 10;
BAR_HEIGHT = 200;


function on_ready()
{
    // just a simple test to ensure works and draws

    draw_read_warm_graph(
        [[1000,1,2,9056.70],
         [5000,1,2,23600.75],
         [10000,1,2,46729.18],
         [50000,1,2,120468.72],
         [100000,1,2,130285.82],
         [150000,1,2,135664.61]]);
}

function NonPerfStats(data_list)
{
    this.num_reads = data_list[NUM_READS_INDEX];
    this.num_threads = data_list[NUM_THREADS_INDEX];
    this.read_type = data_list[READ_TYPE_INDEX];
    this.reads_per_second = data_list[READS_PER_SECOND_INDEX];
}

NonPerfStats.prototype.clone = function()
{
    var data_list_init = [];
    data_list_init[NUM_READS_INDEX] = this.num_reads;
    data_list_init[NUM_THREADS_INDEX] = this.num_threads;
    data_list_init[READ_TYPE_INDEX] = this.read_type;
    data_list_init[READS_PER_SECOND_INDEX] = this.reads_per_second;
    return new NonPerfStats(data_list_init);
};


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
        var condition_data = new NonPerfStats(data_list[index]);
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
    // is a NonPerfStatsObject
    var warm_graph_average_data_list = [];
    
    for (index in warm_graph_data_dict)
    {
        var avg_stats =
            non_perf_stats_avg_throughput(warm_graph_data_dict[index]);
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


/**
 @param {list} stats_array --- Each element is a NonPerfStats object.
 
 @returns {NonPerfStats} --- Contains same num_reads, num_threads, and
 read_type as first object in stats_array.  reads_per_second is
 average value.
 */
function non_perf_stats_avg_throughput(stats_array)
{
    if (stats_array.length == 0)
        throw 'Cannot calulate average value of empty array';
    
    var avg = 0;
    for (var index in stats_array)
        avg += stats_array[index].reads_per_second;
    avg /= stats_array.length;


    var to_return = stats_array[0].clone();
    to_return.reads_per_second = avg;
    return to_return;
}

    