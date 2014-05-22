STACKED_DIV_ID = "stacked_tests";

STACKED_BAR_WIDTH = 80;
STACKED_BAR_SPACING = 40;
STACKED_BAR_HEIGHT = 400;


function stacked_graphs(stacked_graph_data)
{
    var stacked_run_list = process_stacked_data(stacked_graph_data);
    draw_stacked_graphs(stacked_run_list,STACKED_DIV_ID);
}


/**
 @param {list} stacked_run_list --- Each element is a StackedRun
 object.
 */
function draw_stacked_graphs(stacked_run_list,div_id_to_plot_on)
{
    draw_single_stacked_run(stacked_run_list[0],div_id_to_plot_on);
}

function draw_single_stacked_run(stacked_run,div_id_to_plot_on)
{
    var sub_data = stacked_run.averaged_stacked_sub_data;
    var relative_end_time = sub_data.end;
    
    var max_depth = sub_data.max_depth();
    var flattened_data_list = sub_data.flatten(0,max_depth);

    // for (var i in flattened_data_list)
    // {
    //     var flattened_datum = flattened_data_list[i];
    //     flattened_datum.debug_print();
    // }

    var full_graph_width = max_depth * (STACKED_BAR_WIDTH + STACKED_BAR_SPACING);
    
    // indices are strings, values are color strings.  Essentially,
    // there are repeated
    var color_dict = new ColorDict();

    // where to start rectangle x-position form
    var x_rect_positions =
        d3.scale.linear().domain([0, max_depth]).range([20, full_graph_width]);
    // how high to make each rectangle
    var y_heights =
        d3.scale.linear().domain([0,relative_end_time]).
        rangeRound([0,STACKED_BAR_HEIGHT]);

    var bar_chart = d3.select('#' + div_id_to_plot_on).
        append('svg:svg').
        attr('width', full_graph_width).
        attr('height', STACKED_BAR_HEIGHT+60);

    // actually draw the bars
    bar_chart.selectAll('rect').
        data(flattened_data_list).
        enter().
        append('svg:rect').
        attr('x',
             function(flattened_datum, index)
             {
                 return x_rect_positions(flattened_datum.depth);
             }).
        attr('y',
             function(datum)
             {
                 // top left corner
                 var top_left_corner = STACKED_BAR_HEIGHT - y_heights(datum.end);
                 return top_left_corner;
             }).
        attr('height',
             function(datum)
             {
                 return y_heights(datum.end - datum.start);
             }).
        attr('width', STACKED_BAR_WIDTH).
        attr('fill',
             function (datum)
             {
                 return color_dict.get_color(datum.label);
             });

    // add labels to bar regions
    bar_chart.selectAll('text').
        data(flattened_data_list).
        enter().
        append('svg:text').
        attr('x',
             function(flattened_datum)
             {
                 return x_rect_positions(flattened_datum.depth) - 10;
             }).
        attr('y',
             function(flattened_datum)
             {
                 var returned_y =
                     (STACKED_BAR_HEIGHT - y_heights(flattened_datum.end));
                 return returned_y;
             }).
        attr('transform',
             function (flattened_datum)
             {
                 var bar_height = y_heights(flattened_datum.end - flattened_datum.start);
                 var y = (STACKED_BAR_HEIGHT - y_heights(flattened_datum.end) - bar_height);
                 var x = x_rect_positions(flattened_datum.depth) - 10;
                 return 'rotate(270,' + x + ',' + y + ')';
             }).
        text(function(flattened_datum)
             {
                 return flattened_datum.label;
             });

                 
        
    
}

function ColorDict()
{
    this.color_index = 0;
    this.observed_labels = {};
}

ColorDict.prototype.get_color = function (label)
{
    if (!(label in this.observed_labels))
    {
        this.observed_labels[label] = colorbrewer.Set1[9][this.color_index];
        ++this.color_index;
    }
    var color = this.observed_labels[label];
    return color;
};