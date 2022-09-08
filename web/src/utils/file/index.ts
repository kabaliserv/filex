
type UploadOptions = {
    metadata?: Record<string, string>
}

class Upload {
    file: File

    constructor(file: File, option?: UploadOptions) {
        this.file = file
    }

    public start() {

    }
}