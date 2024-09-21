The aim of this project is to rewrite the Tagetik Addin Fast Swapper I have written in C# in Go as an exercise.

Step 1: Swapper as CLI -- done
Step 2: Swapper with TUI -- mostly done


What does it do? 
  The app should let you quickly swap out the Tagetik addin folder (windows default: C:\Tagetik\Tagetik Excel .NET Client)
  with another folder (stored in the same directory) by supplying the other folders name i.e. "Customer1".
  "Customer1" should hold all files needed by the Tagetik Excel Addin for the Tagetik Version of Customer1. 
  The app will then rename the old Tagetik Excel .NET Client folder to a chosen name (i.e. "Customer2") and rename "Customer1"
  to "Tagetik Excel .NET Client".
  The app should also close down MS Excel if it is open and reopen it after swapping. 


To Do:
  -make it so that the swapper can run from anywhere and not only from within the tgk folder (probably a problem with finding settings.json?)
  -add some more info to it >"if you update excel will be closed without further warning" etc
  -Add "new entry"-mode > enter new name (system then should copy and not move the current main folder so that the user can go to excel and set up the new connection etc.)
