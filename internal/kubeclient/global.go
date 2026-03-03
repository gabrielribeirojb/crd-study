package kubeclient

// Isso é um "service locator" simples só pra estudo.
// Em projeto real, seria injeção mais elegante.
var Default ClusterRestoreClient
