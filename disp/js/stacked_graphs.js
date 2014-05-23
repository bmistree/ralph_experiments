STACKED_DIV_ID = "stacked_tests";
STACKED_NOTES_DIV_ID = "stacked_tests_notes";

STACKED_BAR_WIDTH = 40;
STACKED_BAR_SPACING = 10;
STACKED_BAR_ADDITIONAL_SPACING_BETWEEN_EXPERIMENTS = 30;
STACKED_BAR_HEIGHT = 1200;


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
    //draw_single_stacked_run(stacked_run_list[0],div_id_to_plot_on);
    var flattened_sub_data_list = [];
    var max_end_time = 0;
    var max_depth = 0;
    for (var stacked_run_index in stacked_run_list)
    {
        var stacked_run = stacked_run_list[stacked_run_index];
        if (stacked_run.averaged_stacked_sub_data.end > max_end_time)
            max_end_time = stacked_run.averaged_stacked_sub_data.end;

        var depth = stacked_run.averaged_stacked_sub_data.max_depth();
        if (depth > max_depth)
            max_depth = depth;
    }

    // flatten data
    var all_flattened_data = [];
    for (stacked_run_index in stacked_run_list)
    {
        var stacked_run = stacked_run_list[stacked_run_index];
        var additional_fields = {
            num_threads: stacked_run.num_threads,
            graph_index: stacked_run_index
            };
        var run_flattened_data =
            stacked_run.averaged_stacked_sub_data.flatten(
                0,max_depth,additional_fields);
        all_flattened_data = all_flattened_data.concat(run_flattened_data);
    }

    var single_graph_width =
        max_depth * (STACKED_BAR_WIDTH + STACKED_BAR_SPACING);
    var full_graph_width = 
        (single_graph_width + STACKED_BAR_ADDITIONAL_SPACING_BETWEEN_EXPERIMENTS) *
        stacked_run_list.length;
    
    // indices are strings, values are color strings.  Essentially,
    // there are repeated
    var color_dict = new ColorDict();

    // where to start rectangle x-position form
    var x_rect_positions =
        d3.scale.linear().domain([0, max_depth]).range([20, single_graph_width]);
    // how high to make each rectangle
    var y_heights =
        d3.scale.linear().domain([0,max_end_time]).
        rangeRound([0,STACKED_BAR_HEIGHT]);

    var bar_chart = d3.select('#' + div_id_to_plot_on).
        append('svg:svg').
        attr('width', full_graph_width).
        attr('height', STACKED_BAR_HEIGHT);

    // actually draw the bars
    bar_chart.selectAll('rect').
        data(all_flattened_data).
        enter().
        append('svg:rect').
        attr('x',
             function(flattened_datum)
             {
                 var rect_x = rect_x_func(
                     flattened_datum,x_rect_positions,single_graph_width);
                 return rect_x;
             }).
        attr('y',
             function(flattened_datum)
             {
                 var rect_y = rect_y_func(flattened_datum,y_heights);
                 return rect_y;
             }).
        attr('height',
             function(flattened_datum)
             {
                 var rect_height = rect_height_func(flattened_datum,y_heights);
                 return rect_height;
             }).
        attr('width', STACKED_BAR_WIDTH).
        attr('fill',
             function (flattened_datum)
             {
                 return color_dict.get_color(flattened_datum.label);
             }).
        on('mouseover',
           function(flattened_datum)
           {
               var time_ns =
                   flattened_datum.end - flattened_datum.start;
               var time_us = time_ns / 1000;
               $('#' + STACKED_NOTES_DIV_ID).html(
                   'Event: ' + flattened_datum.label + '<br/>' +
                       'Time taken: ' + time_us.toFixed(2) + 'us.');
           }).
        on('mouseout',
           function(flattened_datum)
           {
               $('#' + STACKED_NOTES_DIV_ID).html('<br/><br/>');
           });

    // add labels to bar regions
    bar_chart.selectAll('text').
        data(all_flattened_data).
        enter().
        append('svg:text').
        attr('x',
             function(flattened_datum)
             {
                 return text_x_func(
                     flattened_datum,x_rect_positions,single_graph_width);
             }).
        attr('y',
             function(flattened_datum)
             {
                 return text_y_func(flattened_datum,y_heights);
             }).
        attr('transform',
             function (flattened_datum)
             {
                 // rotate: go to axis specified by x,y and rotate
                 // this element degrees around it.
                 var x = text_x_func(flattened_datum,x_rect_positions,single_graph_width);
                 var y = text_y_func(flattened_datum,y_heights);
                 var transform_to_perform = 'rotate(270,' + x + ',' + y + ')';
                 return transform_to_perform;
             }).
        text(function(flattened_datum)
             {
                 if (flattened_datum.label == 'Total')
                 {
                     var root_time_ns =
                         flattened_datum.end - flattened_datum.start;
                     var root_time_us = root_time_ns / 1000;
                     return 'Num threads: ' + flattened_datum.num_threads +
                         '; Total time: ' + root_time_us.toFixed(2) + 'us';
                 }

                 var rect_height = rect_height_func(flattened_datum,y_heights);
                 if (rect_height < 80)
                     return "";
                 
                 return flattened_datum.label;
             });

}

/**
 Take in a flattened_datum and return the x position that its text
 label should be drawn.
 */
function text_x_func(flattened_datum,x_rect_positions,single_graph_width)
{
    return rect_x_func(
        flattened_datum,x_rect_positions,single_graph_width) + 20;
}
function text_y_func(flattened_datum,y_heights)
{
    return STACKED_BAR_HEIGHT - y_heights(flattened_datum.start) - 10;
}

function rect_x_func(flattened_datum,x_rect_positions,single_graph_width)
{
    var experiment_offset = flattened_datum.graph_index*
        (single_graph_width +
         STACKED_BAR_ADDITIONAL_SPACING_BETWEEN_EXPERIMENTS);
    
    return x_rect_positions(flattened_datum.depth) + experiment_offset;
}

function rect_y_func(flattened_datum,y_heights)
{
    // top left corner
    return STACKED_BAR_HEIGHT - y_heights(flattened_datum.end);
}

function rect_height_func(flattened_datum, y_heights)
{
    return y_heights(flattened_datum.end - flattened_datum.start);
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
        var color_list_to_use = colorbrewer.Set3[12];
        var color_to_use = random_color();
        console.log(color_list_to_use.length);
        if (this.color_index < color_list_to_use.length)
            color_to_use = color_list_to_use[this.color_index];

        this.observed_labels[label] = color_to_use;
        ++this.color_index;
    }
    var color = this.observed_labels[label];
    return color;
};

/**
 From first answer on
 http://stackoverflow.com/questions/1484506/random-color-generator-in-javascript
 */
function random_color()
{
    var letters = '0123456789ABCDEF'.split('');
    var color = '#';
    for (var i = 0; i < 6; i++ )
        color += letters[Math.floor(Math.random() * 16)];
    return color;
}