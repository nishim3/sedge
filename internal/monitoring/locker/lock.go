/*
Copyright 2022 Nethermind

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package locker

//go:generate mockgen -package=mocks -destination=./mocks/locker.go github.com/NethermindEth/sedge/internal/monitoring/locker Locker
type Locker interface {
	// New creates a new Locker instance with the specified path.
	New(path string) Locker

	// Lock acquires the lock, blocking until it is available.
	Lock() error

	// Unlock releases the lock, allowing other processes to acquire it.
	Unlock() error

	// Locked returns true if the lock is currently held, false otherwise.
	Locked() bool
}
