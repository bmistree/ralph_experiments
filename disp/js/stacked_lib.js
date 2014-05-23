(function ()
{
    StackedGraphFactory = new StackedGrapher();    
    function StackedGrapher()
    { }
    StackedGrapher.prototype.params = function (div_id_to_draw_on,div_height)
    {
        return new StackedParams(div_id_to_draw_on,div_height);
    };
    
    function StackedParams(div_id_to_draw_on,div_height)
    {
        this.div_id_to_draw_on = div_id_to_draw_on;
        this.div_height = div_height;
        
        // default parameters
        this.color_dict = new DefaultColorDict();
        this.bar_width = 40;
        this.bar_spacing = 5;
        this.between_experiments_spacing = 30;
        this.hover_func = function(){};
        this.unhover_func = function(){};
    }


    /**
     @param {list} all_data --- Each element is an object that should have the form:

         {
             start: <number>,
             end: <number>,
             label: <string>,
             depth: <int>,
             experiment_index: <int>
         }

     Will draw each rectangle on stacked bar chart relative to its
     start and end position.  label gets displayed on graph next to
     data item (if there is room).  depth indicates what horizontal
     layer of the stacked graph should display rectangle at (left-most
     is zero).  experiment_index is which condition to draw graph in.
     */    
    StackedGrapher.prototype.draw_stacked =
        function(all_data,stacked_params)
    {
        // Find max depth, max end time, and compress data into single
        // list.
        var max_depth = 0;
        var max_end_time  = 0;
        var max_experiment_index = 0;
        for (var i=0; i < all_data.length; ++i)
        {
            var datum = all_data[i];
            if (max_depth < datum.depth)
                max_depth = datum.depth;
            if (max_end_time < datum.end)
                max_end_time = datum.end;
            if (max_experiment_index < datum.experiment_index)
                max_experiment_index = datum.experiment_index;
        }


        var single_graph_width =
            max_depth*(stacked_params.bar_width + stacked_params.bar_spacing);
        var full_graph_width =
            (single_graph_width + stacked_params.between_experiments_spacing)*
            (max_experiment_index + 1);

        // where to draw the top-left of rectangle
        var x_rect_positions =
            d3.scale.linear().domain([0, max_depth]).range([20, single_graph_width]);

        // how high to make each rectangle
        var y_heights =
            d3.scale.linear().domain([0,max_end_time]).
            rangeRound([0,stacked_params.div_height]);


        var bar_chart = d3.select('#' + stacked_params.div_id_to_draw_on).
            append('svg:svg').
            attr('width', full_graph_width).
            attr('height', stacked_params.div_height);


        // draw rectangles
        bar_chart.selectAll('rect').
            data(all_data).
            enter().
            append('svg:rect').
            attr('x',
                 function(datum)
                 {
                     var rect_x = rect_x_func(
                         datum,x_rect_positions,single_graph_width,stacked_params);
                     return rect_x;
                 }).
            attr('y',
                 function(datum)
                 {
                     var rect_y = rect_y_func(datum,y_heights,stacked_params);
                     return rect_y;
                 }).
            attr('height',
                 function(datum)
                 {
                     var rect_height =
                         rect_height_func(datum,y_heights,stacked_params);
                     return rect_height;
                 }).
            attr('width', stacked_params.bar_width).
            attr('fill',
                 function (datum)
                 {
                     return stacked_params.color_dict.get_color(datum.label);
                 }).
            on('mouseover',
               function(datum)
               {
                   stacked_params.hover_func(datum);
               }).
            on('mouseout',
               function(datum)
               {
                   stacked_params.unhover_func(datum);
               });

        // draw labels on bars, if can fit labels
        bar_chart.selectAll('text').
            data(all_data).
            enter().
            append('svg:text').
            attr('x',
                 function(datum)
                 {
                     return text_x_func(
                         datum,x_rect_positions,single_graph_width,stacked_params);
                 }).
            attr('y',
                 function(datum)
                 {
                     return text_y_func(datum,y_heights,stacked_params);
                 }).
            attr('transform',
                 function (datum)
                 {
                     // rotate: go to axis specified by x,y and rotate
                     // this element degrees around it.
                     var x = text_x_func(
                         datum,x_rect_positions,single_graph_width,
                         stacked_params);
                     var y = text_y_func(datum,y_heights,stacked_params);
                     var transform_to_perform = 'rotate(270,' + x + ',' + y + ')';
                     return transform_to_perform;
                 }).
            text(function(datum)
                 {
                     var rect_height = rect_height_func(datum,y_heights);
                     if (rect_height < 80)
                         return "";
                     return datum.label;
                 });
    };
    
    function text_x_func(
        datum,x_rect_positions,single_graph_width,stacked_params)
    {
        return rect_x_func(
            datum,x_rect_positions,single_graph_width,stacked_params) +
            (stacked_params.bar_width/2);
    }
    function text_y_func(datum,y_heights,stacked_params)
    {
        return stacked_params.div_height - y_heights(datum.start) - 10;
    }


    function rect_x_func(datum,x_rect_positions,single_graph_width,stacked_params)
    {
        var experiment_offset = datum.experiment_index*
            (single_graph_width + stacked_params.between_experiments_spacing);
        return x_rect_positions(datum.depth) + experiment_offset;
    }

    function rect_y_func(datum,y_heights,stacked_params)
    {
        // top left corner
        return stacked_params.div_height - y_heights(datum.end);
    }

    function rect_height_func(datum, y_heights)
    {
        return y_heights(datum.end - datum.start);
    }
    
    function DefaultColorDict()
    {
        this.color_index = 0;
        this.observed_labels = {};
        // colors from colorbrewer
        this.color_list_to_use =
            ["#8dd3c7",
             "#ffffb3",
             "#bebada",
             "#fb8072",
             "#80b1d3",
             "#fdb462",
             "#b3de69",
             "#fccde5",
             "#d9d9d9",
             "#bc80bd",
             "#ccebc5",
             "#ffed6f"];
    }


    DefaultColorDict.prototype.get_color = function (label)
    {
        if (!(label in this.observed_labels))
        {
            var color_to_use = random_color();
            if (this.color_index < this.color_list_to_use.length)
                color_to_use = this.color_list_to_use[this.color_index];

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
    
})();
