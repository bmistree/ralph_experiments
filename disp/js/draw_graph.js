

function on_ready()
{
    // just a simple test to ensure works and draws
    draw_read_warm_graph(READ_WARM);
    draw_mem_leak_graph(MEM_LEAK);
    draw_locks_off_graph(LOCKS_OFF);
    wound_wait_graph(WOUND_WAIT);
    gc_off_graph(GC_OFF);
    across_different_ops_graph(PERF_NUM_THREADS);
    stacked_graphs(STACKED_GRAPH_DATA);
    event_map_stacked(EVENT_MAP_STACKED_GRAPH_DATA);
}
