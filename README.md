Kosmos CP1 Tape Tool
====================

With the Kosmos CP1 Tape Tool, you can convert Kosmos CP1 
storage tapes to a binary files and vice versa. 

This allows you to use your old Kosmos CP1 programs
in the Kosmos CP1 emulator. Also, you can write 
new programs in the emulator, and then transfer it
to tape.

To convert a tape, you first need to convert it to a WAV
file. Please make sure that the 16 seconds lead-in is 
completely contained in the file

Usage
-----

### Convert a WAV file to binary:
```kosmos_tape_tool -wav=<wav-file> -bin=<target-binary> read```

### Convert a binary file to WAV:
```kosmos_tape_tool -wav=<target-wav-file> -bin=<tbinary> write```

