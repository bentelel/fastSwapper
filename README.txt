The aim of this project is to rewrite the Tagetik Addin Fast Swapper I have written in C# in Go as an exercise.

Step 1: Swapper as CLI
Step 2: Swapper with GUI


What does it do? 
  Currently, not much. 
  But it is aimed that:
    The app should let you quickly swap out the Tagetik addin folder (windows default: C:\Tagetik\Tagetik Excel .NET Client)
    with another folder (stored in the same directory) by supplying the other folders name i.e. "Customer1".
    "Customer1" should hold all files needed by the Tagetik Excel Addin for the Tagetik Version of Customer1. 
    The app will then rename the old Tagetik Excel .NET Client folder to a chosen name (i.e. "Customer2") and rename "Customer1"
    to "Tagetik Excel .NET Client".
    The app should also close down MS Excel if it is open and reopen it after swapping. 


To Do:
  Test with real data.
  Make it so that excel is killed and restarted on swap.
