// Package runner provides interruptable goroutines.
//
//     task := runner.Go(func(shouldStop runner.S) error {
//       // do setup
//       // defer func(){
//       //   // do teardown
//       // }
//       for {
//         // do some work
//         if shouldStop() {
//           break
//         }
//       }
//       return nil // any errors?
//     })
//
//     // meanwhile...
//     // stop the task
//     task.Stop()
//
//     // wait for it to stop (or time out)
//     select {
//       case <-task.StopChan():
//         // stopped
//       case <-time.After(1 * time.Second):
//         log.Fatalln("task didn't stop quickly enough")
//     }
//
//     // check errors
//     if task.Err() != nil {
//       log.Fatalln("task failed:", task.Err())
//     }
package runner
