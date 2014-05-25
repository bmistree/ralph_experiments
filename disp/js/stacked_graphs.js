STACKED_DIV_ID = 'stacked_tests';
STACKED_NOTES_DIV_ID = 'stacked_tests_notes';
STACKED_SUMMARIZED_JSON_DIV_ID = 'stacked_summarized_json';

STACKED_BAR_WIDTH = 40;
STACKED_BAR_SPACING = 10;
STACKED_BAR_ADDITIONAL_SPACING_BETWEEN_EXPERIMENTS = 30;
STACKED_BAR_HEIGHT = 1200;


function stacked_graphs(stacked_graph_data)
{
    // uncomment if processing raw data
    //var stacked_run_list = process_stacked_data(stacked_graph_data);
    var stacked_run_list =
        deserialize_serialized_stacked_run_list(stacked_graph_data);
    draw_stacked_graphs(stacked_run_list,STACKED_DIV_ID, STACKED_NOTES_DIV_ID);
    update_summarized_json(stacked_run_list,STACKED_SUMMARIZED_JSON_DIV_ID);
}

/**
 @param {list} --- Each element should be an object generated from
 calling summarized_json on a StackedRun object.

 @return {list} --- Each element is a StackedRun object.
 */
function deserialize_serialized_stacked_run_list(serialized_stacked_run_list)
{
    var to_return = [];
    for(var i = 0; i < serialized_stacked_run_list.length; ++i)
    {
        var serialized_stacked_run = serialized_stacked_run_list[i];
        to_return.push(
            deserialize_serialized_stacked_run(serialized_stacked_run));
    }
    return to_return;
}

/**
 @param {Object} serialized_stacked_run --- Object generated from
 importing json generated from StackedRun's summarized_json call.

 @return {StackedRun}
 */
function deserialize_serialized_stacked_run(serialized_stacked_run)
{
    var serialized_sub_data = serialized_stacked_run.averaged_stacked_sub_data;
    var throughput = serialized_stacked_run.throughput;
    var num_threads = serialized_stacked_run.num_threads;

    var sub_data = deserialize_serialized_sub_data(serialized_sub_data);
    return new StackedRun(throughput,num_threads,sub_data);
}

/**
 @param {Object} serialized_sub_data --- Generated from importing json
 generated from calling summarized_json on a StackedSubData object.

 @param {StackedSubData}
 */
function deserialize_serialized_sub_data(serialized_sub_data)
{
    var time = serialized_sub_data.time;
    var label = serialized_sub_data.label;
    var start = serialized_sub_data.start;
    var to_return = new StackedSubData(time,label,start);

    var children = serialized_sub_data.children;
    for (var i = 0; i < children.length; ++i)
    {
        var child = children[i];
        var child_as_sub_data = deserialize_serialized_sub_data(child);
        to_return.add_child(child_as_sub_data);
    }
    return to_return;
}


function update_summarized_json(stacked_run_list,div_id_to_update_on)
{
    var json_summary_string = '[';
    for (var stacked_run_index in stacked_run_list)
    {
        var stacked_run = stacked_run_list[stacked_run_index];
        json_summary_string += stacked_run.summarized_json();
        if (stacked_run_index != (stacked_run_list.length - 1))
            json_summary_string += ',';
    }
    json_summary_string += ']';
    $('#' + div_id_to_update_on).html(json_summary_string);
}


/**
 @param {list} stacked_run_list --- Each element is a StackedRun
 object.
 */
function draw_stacked_graphs(stacked_run_list,div_id_to_plot_on,notes_div_id)
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
            experiment_index: stacked_run_index
            };
        var run_flattened_data =
            stacked_run.averaged_stacked_sub_data.flatten(
                0,max_depth,additional_fields);
        all_flattened_data = all_flattened_data.concat(run_flattened_data);
    }

    var params = StackedGraphFactory.params(div_id_to_plot_on,STACKED_BAR_HEIGHT);
    params.hover_func = function (datum)
    {
       var time_ns =
           datum.end - datum.start;
       var time_us = time_ns / 1000;
       $('#' + notes_div_id).html(
           'Event: ' + datum.label + '<br/>' +
               'Time taken: ' + time_us.toFixed(2) + 'us.');
    };
    params.unhover_func = function(datum)
    {
        $('#' + notes_div_id).html('<br/><br/>');
    };
    StackedGraphFactory.draw_stacked(all_flattened_data,params);
}