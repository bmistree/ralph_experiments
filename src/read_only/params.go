package read_only

type Parameter struct {
    numReads uint32
    numThreads uint32
    opType operationType
    persistentThreadPoolSize uint32
    maxThreadPoolSize uint32
    atomIntUUIDGeneration bool
    woundWaitOn bool
    gcOn bool
    readsOnOtherAtomNum bool
    perfOn bool
}

func createDefaultParameter () *Parameter {
    toReturn := Parameter {
        numReads: 100000,
        numThreads: 1,
        opType: READ_ATOM_NUM,
        persistentThreadPoolSize: 0,
        maxThreadPoolSize: 0,
        atomIntUUIDGeneration: false,
        woundWaitOn: false,
        gcOn: true,
        readsOnOtherAtomNum: false,
        perfOn: false,
    }
    return &toReturn
}